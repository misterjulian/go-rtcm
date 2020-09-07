
RTCM parser library implementation in Go.

```
Usage example:

import rtcm "github.com/misterjulian/go-rtcm"

...
conn, err := d.DialContext(ctx, "tcp", "some address here")
...

source = rtcm.New(conn)
for {
    frame, err := source.NextFrame() 
    ...
    Your logic here
    ...
}
```
