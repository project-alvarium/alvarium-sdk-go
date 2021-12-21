/*******************************************************************************
 * Copyright 2021 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/
package interfaces

import (
	"context"
	"sync"
)

type Sdk interface {
	// BootstrapHandler provides a hook whereby shutdown signals can be trapped and gracefull handled in order to clean
	// up resources and active connections.
	BootstrapHandler(ctx context.Context, wg *sync.WaitGroup) bool

	// Create handles annotations relative to the creation of new data
	// The data parameter is the given piece of data being created marshalled as a byte array. This in turn will be
	// used to generate a hash uniquely identifying the data item.
	Create(ctx context.Context, data []byte)

	// Mutate handles annotations relative to a data modification. That is to say, an older piece of data is being
	// updated or transformed into new data.
	// The old, new parameters are the given data elements marshalled as byte arrays. You must have byte representations
	// of both the old and the new data in order to establish a provenance linkage through Mutate.
	Mutate(ctx context.Context, old, new []byte)

	// Transit is a proposed method for cases where an existing piece of data is received by a separate application
	// that is not the originator of the data. The data has simply transited from one application/host to the other.
	// This method could be used to asses the signature validity on the received data, secure comms, checksum validation,
	// etc.
	Transit(ctx context.Context, data []byte)

	// Publish is proposed to provide extensibility for annotators that may need to attest to the state of data before it
	// is sent over the wire. Publish could also be useful in cases where the downstream host receiving the data isn't
	// running Alvarium-enabled applications.
	Publish(ctx context.Context, data []byte)
}
