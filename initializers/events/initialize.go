// SPDX-FileCopyrightText: 2025 Nimble Tech
// SPDX-License-Identifier: MIT

package events

import "github.com/gookit/event"

func InitializeManager() *event.Manager {
	return event.NewManager("main")
}
