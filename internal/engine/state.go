package engine

import "poker_tui/internal/game"

// Screen represents the current TUI state/screen.
type Screen int

const (
	ScreenMainIdle        Screen = iota // No hand; ready to bet/deal
	ScreenHandDealt                     // 5 cards shown; hold phase
	ScreenHandResolved                  // Result shown; gamble eligible
	ScreenGambleStage                   // Active gamble stage
	ScreenGambleResult                  // Gamble resolved (win/lose)
	ScreenOptionsMenu                   // Options overlay
	ScreenHelpOverlay                   // Help + paytable overlay
	ScreenHighScoreScreen               // Stats screen
	ScreenErrorScreen                   // Lock/error screen
)

func (s Screen) String() string {
	return [...]string{
		"MainIdle", "HandDealt", "HandResolved",
		"GambleStage", "GambleResult",
		"OptionsMenu", "HelpOverlay", "HighScoreScreen", "ErrorScreen",
	}[s]
}

// CardDesign controls the card rendering style.
type CardDesign int

const (
	DesignClassic CardDesign = iota
	DesignMinimal
	DesignWide
)

func (d CardDesign) String() string {
	return [...]string{"classic", "minimal", "wide"}[d]
}

// Theme is a lipgloss colour preset.
type Theme int

const (
	ThemeDark Theme = iota
	ThemeAmber
	ThemeGreen
	ThemeMono
)

func (t Theme) String() string {
	return [...]string{"dark", "amber", "green", "mono"}[t]
}

// GambleStep records one stage of the gamble mini-game.
type GambleStep struct {
	Card      game.Card
	Choice    string // "red" | "black"
	Outcome   string // "win" | "lose"
	PotBefore int
	PotAfter  int
}

// GambleState holds all live gamble-round data.
type GambleState struct {
	Active     bool
	Stage      int          // stages completed (0 = none yet)
	MaxStages  int          // default 5
	CurrentPot int          // credits currently at risk
	History    []GambleStep // completed steps
	// Awaiting the player's choice for this stage:
	CurrentCard game.Card // drawn but not yet resolved
	Revealed    bool      // true once the card is shown
}

// Stats tracks cumulative gameplay statistics.
type Stats struct {
	HandsPlayed   int
	HandsWon      int
	HandsLost     int
	BiggestWin    int // largest single payout
	TotalWon      int // sum of all payouts
	TotalWagered  int // sum of all bets
	CurrentStreak int // positive=win streak, negative=loss streak
	BestStreak    int
	LifetimeDelta int // credits gained vs starting 100
	GambleWins    int
	GambleLosses  int
}

// Options holds player-configurable settings.
type Options struct {
	CardDesign CardDesign
	Theme      Theme
	AutoHold   bool // automatically hold high cards (J+) and pairs after deal
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		CardDesign: DesignClassic,
		Theme:      ThemeDark,
		AutoHold:   false,
	}
}

// GameState is the single authoritative model for the entire game.
type GameState struct {
	SaveVersion int    // bump when breaking changes to save format
	Screen      Screen // current screen/state
	ReadOnly    bool   // true when running as second instance

	// Economy
	Credits int
	Bet     int

	// XP / level
	XP    int
	Level int

	// Current hand
	Hand       game.Hand
	Deck       game.Deck
	LastResult game.HandResult

	// Gamble
	Gamble GambleState

	// Statistics
	Stats Stats

	// Settings
	Options Options

	// UI transient messages (not persisted to save)
	Message      string
	PrevScreen   Screen // for back-navigation from overlays
	ErrorMessage string
}

const (
	SaveVersion     = 1
	DefaultCredits  = 100
	DefaultBet      = 5
	MinBet          = 1
	MaxBet          = 50
	MaxGambleStages = 5
)

// NewGameState creates a brand-new game state at defaults.
func NewGameState() *GameState {
	return &GameState{
		SaveVersion: SaveVersion,
		Screen:      ScreenMainIdle,
		Credits:     DefaultCredits,
		Bet:         DefaultBet,
		Options:     DefaultOptions(),
		Gamble:      GambleState{MaxStages: MaxGambleStages},
	}
}
