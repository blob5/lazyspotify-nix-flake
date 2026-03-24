package logger

import (
	"github.com/dubeyKartikay/lazyspotify/core/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var Log zerolog.Logger

func init(){
	lumberjackLogger := utils.NewLumberjackLogger("lazyspotify.log")
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	Log = zerolog.New(lumberjackLogger).With().Timestamp().Caller().Stack().Logger()
	Log.Debug().Msg("logger initialized")
	Log.Debug().Msgf("config directory: %s", utils.SafeGetConfigDir())

}
