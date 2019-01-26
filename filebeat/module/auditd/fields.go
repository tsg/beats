// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// Code generated by beats/dev-tools/cmd/asset/asset.go - DO NOT EDIT.

package auditd

import (
	"github.com/elastic/beats/libbeat/asset"
)

func init() {
	if err := asset.SetFields("filebeat", "auditd", asset.ModuleFieldsPri, AssetAuditd); err != nil {
		panic(err)
	}
}

// AssetAuditd returns asset data.
// This is the base64 encoded gzipped contents of module/auditd.
func AssetAuditd() string {
	return "eJy8VsGO2zYQvfsrHnJKD2vk7AIFimxTGGibou3docmxNFiao5LDddWvL0jZXkWWNzWyG95Eke+9eTMc8g4P1K9gsmN1C0BZPa3w5sc68WYBOEo2cqcsYYUfFgDwq7jsCTuJ6ExMHJojALw0abkAUitRN1bCjpsVNGZaADsm79KqQtwhmD2NiMvQvqMVmii5O87MkJfxoSJhF2UPbWnKXsaYbEzopTnPzTE+wzrH/AuH/M/AX6CX+E0UxvsjP0wkOLF5T0HJoaVI2JI1OdFnuNpSPyzug9mzhQkOjyb22PZHeHqkoFXxcrR1Guc41khWotuULZ/9fzbCMv46eXqF9IlCvNuYzO42/A8Si1scBvwEbTmBU7VUvDuSr++RE7laaOVPThTRRS6fcgFaMSrqvNRAh5eWGuhwlrocuba+hzUBWxr0z4mNxhJ2WXOkM7I8RVmmwLsywRHsKChrD9ua0FDCW88P05yiFJbsy2mMIvrd9YQlSi+cr0QpsYRXyNjLai0Je9K6xFoniYIywVyA1jhUJgnb9mOw2RAS/Z0p2Gmuhr7jJTRfdzBP8Ah5v6U4r8FYq7fT1JCNtZKDViCYlMSyKX3swNpWQ6uMedru1qNWWNf3kFr16KJYSukK9muCs9L+xqIr6EMKCkMFAAeY8JxBZdntNGXX0PXBwbE1SgmHlm07/JKspyhV1PipruUMKqfTKk74l6LcbU0i9z0MHo3PVDa/w55MSGA9HaYdx6QV9Erdvbs9ugHTxKbemKeWmPpUQrPG+3mqeGuXKFyRUvZns0YceJuyLfUBidgZ9jnSlXbakHA3e7inj4r/Ieq9BDUcEn4mWf8ODjuJe1NWozFa3g4ONTGQUDV/Gh4+y/LyMM7FTxeQtVCW+Bh8jy5SKqYO18qR5CdvkrJNZKJt0flcOicnmEfD3mw9lZfIXD9000Kae4mMrbISlAMF3ZTvi2Un3x6oP0i8pPyiezidw9Krjkk9c15W/VkWa/8tFbH218VEaljCN5QzED7jTmn/sd9wko0V90qq3g8sWP/5EYXluh4vth6Iqzoakk0nHKY33g3+lGuZNbta+vBG68cXk/a6Fv1RSUYO/RcAAP//EHqm1w=="
}