// Copyright 2018-2025 The Olric Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package olric

import (
	"errors"

	"github.com/olric-data/olric/internal/protocol"
	"github.com/olric-data/olric/internal/server"
	"github.com/tidwall/redcon"
)

// authCommandHandler handles authentication requests sent by clients and verifies the provided password for access.
func (db *Olric) authCommandHandler(conn redcon.Conn, cmd redcon.Command) {
	authCmd, err := protocol.ParseAuthCommand(cmd)
	if err != nil {
		protocol.WriteError(conn, err)
		return
	}

	if !db.config.Authentication.Enabled() {
		protocol.WriteError(conn, errors.New("AUTH <password> called without any password configured for the default user. Are you sure your configuration is correct?"))
		return
	}

	if authCmd.Password == db.config.Authentication.Password {
		ctx := conn.Context().(*server.ConnContext)
		ctx.SetAuthenticated(true)
		conn.WriteString(protocol.StatusOK)
		return
	}
	protocol.WriteError(conn, ErrWrongPass)
}
