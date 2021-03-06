# listenerutils
`import "github.com/carousell/Orion/utils/listenerutils"`

* [Overview](#pkg-overview)
* [Imported Packages](#pkg-imports)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>

## <a name="pkg-imports">Imported Packages</a>

- [github.com/carousell/Orion/utils/log](./../log)

## <a name="pkg-index">Index</a>
* [type CustomListener](#CustomListener)
  * [func NewListener(network, laddr string) (CustomListener, error)](#NewListener)
  * [func NewListenerWithTimeout(network, laddr string, timeout time.Duration) (CustomListener, error)](#NewListenerWithTimeout)

#### <a name="pkg-files">Package files</a>
[listenerutils.go](./listenerutils.go) 

## <a name="CustomListener">type</a> [CustomListener](./listenerutils.go#L14-L19)
``` go
type CustomListener interface {
    net.Listener
    CanClose(bool)
    GetListener() CustomListener
    StopAccept()
}
```
CustomListener provides an implementation for a custom net.Listener

### <a name="NewListener">func</a> [NewListener](./listenerutils.go#L155)
``` go
func NewListener(network, laddr string) (CustomListener, error)
```
NewListener creates a new CustomListener

### <a name="NewListenerWithTimeout">func</a> [NewListenerWithTimeout](./listenerutils.go#L159)
``` go
func NewListenerWithTimeout(network, laddr string, timeout time.Duration) (CustomListener, error)
```

- - -
Generated by [godoc2ghmd](https://github.com/GandalfUK/godoc2ghmd)