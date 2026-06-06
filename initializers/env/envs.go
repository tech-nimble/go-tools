// SPDX-FileCopyrightText: 2025 Nimble Tech
// SPDX-License-Identifier: MIT

package env

import (
	"github.com/gobuffalo/envy"
	"github.com/rs/zerolog/log"
)

func InitializeEnv() {
	if err := envy.Load(); err != nil {
		log.Info().Err(err).Msg("can not load .env file")
	}
}
