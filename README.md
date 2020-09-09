
RTCM parser library implementation in Go.

```
Usage example:

import rtcm "github.com/misterjulian/go-rtcm"

...
conn, err := d.DialContext(ctx, "tcp", "some address here")
...

scanner = rtcm.NewScanner(conn)
for {
    frame, err := scanner.ScanFrame() 
    ...
    Your logic here
    ...
}
```
