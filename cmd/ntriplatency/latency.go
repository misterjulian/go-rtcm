package main

import (
    "github.com/geoscienceaustralia/go-rtcm/pkgs/rtcm3"
    "github.com/umeat/go-ntrip/pkgs/ntrip"
    "fmt"
    "os"
    "time"
)

func main() {
    r, _ := ntrip.Client(os.Args[1], os.Args[2], os.Args[3])
    scanner := rtcm3.NewScanner(r)
    for msg, err := scanner.Next(); err == nil; msg, err = scanner.Next() {
        switch int(msg.Number()) {
            case 1071, 1072, 1073, 1074, 1075, 1076, 1077,
                 1081, 1082, 1083, 1084, 1085, 1086, 1087,
                 1091, 1092, 1093, 1094, 1095, 1096, 1097,
                 1101, 1102, 1103, 1104, 1105, 1106, 1107,
                 1111, 1112, 1113, 1114, 1115, 1116, 1117,
                 1121, 1122, 1123, 1124, 1125, 1126, 1127,
                 1001, 1002, 1003, 1004, 1009, 1010, 1011, 1012:

                fmt.Println(msg.Number(), time.Now().UTC().Sub(msg.(rtcm3.Observable).Time()))

            default:
                fmt.Println(msg.Number())
        }
    }
}
