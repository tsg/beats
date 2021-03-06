// Package checkpoint persists event log state information to disk so that
// event log monitoring can resume from the last read event in the case of a
// restart or unexpected interruption.
package checkpoint

import (
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"io/ioutil"

	"github.com/elastic/beats/libbeat/logp"
	"gopkg.in/yaml.v2"
)

// Checkpoint persists event log state information to disk.
type Checkpoint struct {
	wg            sync.WaitGroup // WaitGroup used to wait on the shutdown of the checkpoint worker.
	done          chan struct{}  // Channel for shutting down the checkpoint worker.
	once          sync.Once      // Used to guarantee shutdown happens once.
	file          string         // File where the state is persisted.
	numUpdates    int            // Number of updates received since last persisting to disk.
	maxUpdates    int            // Maximum number of updates to buffer before persisting to disk.
	flushInterval time.Duration  // Maximum time interval that can pass before persisting to disk.
	sort          []string       // Slice used for sorting states map (store to save on mallocs).

	lock   sync.RWMutex
	states map[string]EventLogState

	save chan EventLogState
}

// PersistedState represents the format of the data persisted to disk.
type PersistedState struct {
	UpdateTime time.Time       `yaml:"update_time"`
	States     []EventLogState `yaml:"event_logs"`
}

// EventLogState represents the state of an individual event log.
type EventLogState struct {
	Name         string    `yaml:"name"`
	RecordNumber uint32    `yaml:"record_number"`
	Timestamp    time.Time `yaml:"timestamp"`
}

// NewCheckpoint creates and returns a new Checkpoint. This method loads state
// information from disk if it exists and starts a goroutine for persisting
// state information to disk. Shutdown should be called when finished to
// guarantee any in-memory state information is flushed to disk.
//
// file is the name of the file where event log state is persisted as YAML.
// maxUpdates is the maximum number of updates checkpoint will accept before
// triggering a flush to disk. interval is maximum amount of time that can
// pass since the last flush before triggering a flush to disk (minimum value
// is 1s).
func NewCheckpoint(file string, maxUpdates int, interval time.Duration) (*Checkpoint, error) {
	c := &Checkpoint{
		done:          make(chan struct{}),
		file:          file,
		maxUpdates:    maxUpdates,
		flushInterval: interval,
		sort:          make([]string, 0, 10),
		states:        make(map[string]EventLogState),
		save:          make(chan EventLogState, 1),
	}

	// Minimum batch size.
	if c.maxUpdates < 1 {
		c.maxUpdates = 1
	}

	// Minimum flush interval.
	if c.flushInterval < time.Second {
		c.flushInterval = time.Second
	}

	// Read existing state information:
	ps, err := c.read()
	if err != nil {
		return nil, err
	}

	if ps != nil {
		for _, state := range ps.States {
			c.states[state.Name] = state
		}
	}

	c.wg.Add(1)
	go c.run()
	return c, nil
}

// run is worker loop that reads incoming state information from the save
// channel and persists it when the number of changes reaches maxEvents or
// the amount of time since the last disk write reaches flushInterval.
func (c *Checkpoint) run() {
	defer c.wg.Done()

	flushTimer := time.NewTimer(c.flushInterval)
	defer flushTimer.Stop()
loop:
	for {
		select {
		case <-c.done:
			break loop
		case s := <-c.save:
			c.lock.Lock()
			c.states[s.Name] = s
			c.lock.Unlock()
			c.numUpdates++
			if c.numUpdates < c.maxUpdates {
				continue
			}
		case <-flushTimer.C:
		}

		c.persist()
		flushTimer.Reset(c.flushInterval)
	}

	c.persist()
}

// Shutdown stops the checkpoint worker (which persists any state to disk as
// it stops). This method blocks until the checkpoint worker shutdowns. Calling
// this method more once is safe and has no effect.
func (c *Checkpoint) Shutdown() {
	c.once.Do(func() {
		close(c.done)
		c.wg.Wait()
	})
}

// States returns the current in-memory event log state. This state information
// is bootstrapped with any data found on disk at creation time.
func (c *Checkpoint) States() map[string]EventLogState {
	c.lock.RLock()
	defer c.lock.RUnlock()

	copy := make(map[string]EventLogState)
	for k, v := range c.states {
		copy[k] = v
	}

	return copy
}

// Persist queues the given event log state information to be written to disk.
func (c *Checkpoint) Persist(name string, recordNumber uint32, ts time.Time) {
	c.save <- EventLogState{
		Name:         name,
		RecordNumber: recordNumber,
		Timestamp:    ts,
	}
}

// persist writes the current state to disk if the in-memory state is dirty.
func (c *Checkpoint) persist() bool {
	if c.numUpdates == 0 {
		return false
	}

	err := c.flush()
	if err != nil {
		logp.Err("%v", err)
		return false
	}

	logp.Debug("checkpoint", "Checkpoint saved to disk. numUpdates=%d",
		c.numUpdates)
	c.numUpdates = 0
	return true
}

// flush writes the current state to disk.
func (c *Checkpoint) flush() error {
	tempFile := c.file + ".new"
	file, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("Failed to flush state to disk. Could not open %s. %v",
			tempFile, err)
	}

	// Sort persisted eventLogs by name.
	c.sort = c.sort[:0]
	for k := range c.states {
		c.sort = append(c.sort, k)
	}
	sort.Strings(c.sort)

	ps := PersistedState{
		UpdateTime: time.Now().UTC(),
		States:     make([]EventLogState, len(c.sort)),
	}
	for i, name := range c.sort {
		ps.States[i] = c.states[name]
	}

	data, err := yaml.Marshal(ps)
	if err != nil {
		file.Close()
		return fmt.Errorf("Failed to flush state to disk. Could not marshal "+
			"data to YAML. %v", err)
	}

	_, err = file.Write(data)
	if err != nil {
		file.Close()
		return fmt.Errorf("Failed to flush state to disk. Could not write to "+
			"%s. %v", tempFile, err)
	}

	file.Close()
	err = os.Rename(tempFile, c.file)
	return err
}

// read loads the persisted state from disk. If the file does not exists then
// the method returns nil and no error.
func (c *Checkpoint) read() (*PersistedState, error) {
	contents, err := ioutil.ReadFile(c.file)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return nil, err
	}

	ps := &PersistedState{}
	err = yaml.Unmarshal(contents, ps)
	if err != nil {
		return nil, err
	}

	return ps, nil
}
