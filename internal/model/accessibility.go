package model

type AccessibilityProfile struct {
	UserID uint64 `gorm:"primaryKey" json:"user_id"`
	User   User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	VisionImpaired    bool `json:"tuna_netra"`
	HearingImpaired   bool `json:"tuna_rungu"`
	PhysicalImpaired  bool `json:"tuna_daksa"`
	CognitiveImpaired bool `json:"kesulitan_kognitif"` // atau tuna_grahita
	SpeechImpaired    bool `json:"tuna_wicara"`

	ScreenReaderCompatible bool `json:"kompatibel_screen_reader"`
	AudioDescription       bool `json:"deskripsi_audio"`
	SubtitlesRequired      bool `json:"butuh_subtitle"`
	VisualNotifications    bool `json:"notifikasi_visual"`
	KeyboardNavigation     bool `json:"navigasi_keyboard"`
	VoiceCommand           bool `json:"perintah_suara"`
	AISummary              bool `json:"ringkasan_ai"`
	FocusMode              bool `json:"mode_fokus"`
	TextBasedSubmission    bool `json:"pengumpulan_teks"`
}
