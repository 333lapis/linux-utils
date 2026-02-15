# playerctl-monitor

Outputs information from playerctl at a regular interval in a user-specified format. For my Swaybar config but can probably be used in other contexts.

## Building

`go build -o playerctl-monitor main.go` or put it somewhere else like in your path

## Usage

see `playerctl-monitor -h`

## Formatting

Use `-f` to specify the format of the output, using the following placeholders

| placeholder | name |
| - | - |
| `@t@` | title |
| `@a@` | artist |
| `@A@` | album artist |
| `@al@` | album |
| `@au@` | album art URL |
| `@l@` | song length (in microseconds) |
| `@lF@` | song length formatted as M:SS |
| `@s@` | playback status (true if playing, false if paused) |
| `@p@` | position in song (in seconds, floating-point) |
| `@pF@` | position formatted as M:SS |
| `@v@` | volume |
| `@L@` | current loop setting, can be `None`, `Playlist`, or `Track` |
| `@S@` | current shuffle setting (true for on, false for off) |

## Examples

```sh
playerctl-monitor -f "@a@ - @t@    @pF@"
# 2hollis - cliche    2:21

playerctl-monitor -f "[ @pF@ / @lF@ ]    @t@ on @al@ by @a@"
# [ 0:04 / 2:28 ]    Stateside + Bladee on Fancy Some More? by PinkPantheress
```