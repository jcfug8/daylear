package plugins

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// NewSnowflake -
func NewSnowflake(log zerolog.Logger) *Snowflake {
	return &Snowflake{log: log}
}

// Snowflake -
type Snowflake struct {
	log zerolog.Logger
}

func (s *Snowflake) Name() string {
	return "Snowflake"
}

func getZoneId() int {
	if os.Getenv("ZONE_ID") == "" {
		return 0
	}

	id, err := strconv.Atoi(os.Getenv("ZONE_ID"))
	if err != nil {
		panic(fmt.Errorf("invalid ZONE_ID: %v", err))
	}

	return id
}

func (s *Snowflake) Initialize(db *gorm.DB) error {
	var count int64

	err := db.Table("pg_proc").Where("proname = ?", "id_generator").Count(&count).Error
	if err != nil {
		return err
	} else if count > 0 {
		return nil
	}

	s.log.Info().Msg("creating id_generator function")

	tx := db.Exec(fmt.Sprintf(`
		-- sequence used by the id generator that avoids collisions in the same millisecond
		CREATE SEQUENCE IF NOT EXISTS public.global_id_seq;
		ALTER SEQUENCE IF EXISTS public.global_id_seq CACHE 10000;
	
		-- this generates a snowflake id-like sequence
		-- 43 bits of timestamp in milliseconds
		-- 8 bits of zone id (zone is equivalent to a cluster, part of a region)
		-- 12 bits of sequence to avoid collisions for ids generated in the same millisecond
	
		CREATE OR REPLACE FUNCTION public.id_generator()
			RETURNS bigint
			LANGUAGE 'plpgsql'
		AS $BODY$
		DECLARE
			seq_id bigint;
			now_millis bigint;
			-- this must be unique for each zone, globally. 8 bit so 0-255
			-- suggested usage is 10 * region_id + cluster_id (e.g. 10-19 for "usa", 20-29 for "eur")
			zone_id int := %d; -- dev can use zone ids from 0 to 9
			result bigint:= 0;
		BEGIN
			SELECT nextval('public.global_id_seq') %% 4096 INTO seq_id;
	
			SELECT FLOOR(EXTRACT(EPOCH FROM clock_timestamp()) * 1000) INTO now_millis;
			result := (now_millis) << 20;
			result := result | (zone_id << 12);
			result := result | (seq_id);
			return result;
		END;
		$BODY$;
	`, getZoneId()))

	return tx.Error
}
