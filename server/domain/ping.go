package domain

func (d *Domain) Ping() error {
	d.log.Debug().Msg("Pinging domain")
	return nil
}
