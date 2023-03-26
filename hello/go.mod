module prova/hello

        go 1.20

        replace prova/greetings => ../greetings

        replace prova/log_utils => ../log_utils

        require prova/greetings v0.0.0-00010101000000-000000000000

        require (
        github.com/rs/xid v1.4.0
        prova/log_utils v0.0.0-00010101000000-000000000000
        )

        require (
        github.com/sirupsen/logrus v1.9.0 // indirect
        golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
        )
