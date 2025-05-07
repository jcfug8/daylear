# postgres

Provides postgres dialects for use in the gorm dialer.

## Usage

```golang
import (
    gormdialer "github.com/jcfug8/daylear/server/adapters/clients/gorm/dialer"
    "github.com/jcfug8/daylear/server/adapters/clients/gorm/dialer/dialects/postgres"
)

func main() {
    fx.New(
        gormdialer.Module,
        postgres.Module,
    ).Run()
}
```

### Multiple Dialects

If you need to use multiple dialects, you can use the `ConfigsTag` and the
`DialectsTag` to provide and consume the postgres configs and postgres dialects
respectively.
