# monitor

Outputs information at a regular interval in a user-specified format. For my Swaybar config but can probably be used in other contexts.

## Building

`go build -o $HOME/.config/sway/monitor main.go` for example

## Usage

see `monitor -h`

## Formatting

Use `-f` to specify the format of the output, using the following placeholders

### PulseAudio

| placeholder | name |
| - | - |
| `@pa:m@` | whether PulseAudio is muted (boolean) |
| `@pa:v@` | volume of front-left output |

### Time

| placeholder | name |
| - | - |
| `@t:h@` | current hour |
| `@t:m@` | current minute |
| `@t:s@` | current second |

### Battery

| placeholder | name |
| - | - |
| `@b:c@` | capacity (percentage) |

### playerctl

| placeholder | name |
| - | - |
| `@p:t@` | title |
| `@p:a@` | artist |
| `@p:A@` | album artist |
| `@p:al@` | album |
| `@p:au@` | album art URL |
| `@p:l@` | song length (in microseconds) |
| `@p:lF@` | song length formatted as M:SS |
| `@p:s@` | playback status (true if playing, false if paused) |
| `@p:p@` | position in song (in seconds, floating-point) |
| `@p:pF@` | position formatted as M:SS |
| `@p:v@` | volume |
| `@p:L@` | current loop setting, can be `None`, `Playlist`, or `Track` |
| `@p:S@` | current shuffle setting (true for on, false for off) |

## Examples

todo add examples