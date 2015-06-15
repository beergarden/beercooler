# beercooler
ビールを冷やします

## Usage

Compile and install:

```
go install
```

Get current temperature and turn the fan ON/OFF:

```
beercooler {GPIO for fan controller} {upper temperature limit}
```

Execute periodically with crontab:

```
*/5 * * * * /home/pi/go/bin/beercooler {GPIO for fan controller} {upper temperature limit}
```

