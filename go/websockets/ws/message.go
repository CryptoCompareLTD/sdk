package ws

// CADLITickMsg describes a CADLI Tick received from the CCData streamer
type CADLITickMsg struct {
	MsgType                       string  `json:"TYPE"`
	Instrument                    string  `json:"INSTRUMENT"`
	Market                        string  `json:"MARKET"`
	ValueFlag                     string  `json:"VALUE_FLAG"`
	CCSEQ                         uint64  `json:"CCSEQ"`
	Value                         float64 `json:"VALUE"`
	ValueLastUpdateTS             uint64  `json:"VALUE_LAST_UPDATE_TS"`
	ValueLastUpdateTSNS           uint64  `json:"VALUE_LAST_UPDATE_TS_NS"`
	CurrHourVol                   float64 `json:"CURRENT_HOUR_VOLUME"`
	CurrHourQuoteVol              float64 `json:"CURRENT_HOUR_QUOTE_VOLUME"`
	CurrHourVolTopTier            float64 `json:"CURRENT_HOUR_VOLUME_TOP_TIER"`
	CurrHourQuoteVolTopTier       float64 `json:"CURRENT_HOUR_QUOTE_VOLUME_TOP_TIER"`
	CurrHourVolDirect             float64 `json:"CURRENT_HOUR_VOLUME_DIRECT"`
	CurrHourQuoteVolDirect        float64 `json:"CURRENT_HOUR_QUOTE_VOLUME_DIRECT"`
	CurrHourVolTopTierDirect      float64 `json:"CURRENT_HOUR_VOLUME_TOP_TIER_DIRECT"`
	CurrHourQuoteVolTopTierDirect float64 `json:"CURRENT_HOUR_QUOTE_VOLUME_TOP_TIER_DIRECT"`
	CurrHourOpen                  float64 `json:"CURRENT_HOUR_OPEN"`
	CurrHourHigh                  float64 `json:"CURRENT_HOUR_HIGH"`
	CurrHourLow                   float64 `json:"CURRENT_HOUR_LOW"`
	CurrHourTotalIndexUpdates     uint64  `json:"CURRENT_HOUR_TOTAL_INDEX_UPDATES"`
	CurrHourChange                float64 `json:"CURRENT_HOUR_CHANGE"`
	CurrHourChangePerc            float64 `json:"CURRENT_HOUR_CHANGE_PERCENTAGE"`
}
