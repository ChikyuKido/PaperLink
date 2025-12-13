package entity

var TagColors = []string{
	"#E27D60", "#85DCB", "#E8A87C", "#C38D9E", "#41B3A3",
	"#6B5B95", "#88B04B", "#F7CAC9", "#92A8D1", "#034F84",
	"#F7786B", "#79C753", "#B565A7", "#955251", "#009B77",
	"#DD4124", "#D65076", "#45B8AC", "#98B4D4", "#F4A688",
	"#C3447A", "#6A5ACD", "#20B2AA", "#FF7F50", "#5F9EA0",
	"#FFD700", "#DA70D6", "#7B68EE", "#CD5C5C", "#4682B4",
	"#9ACD32", "#FF69B4", "#6495ED", "#40E0D0", "#FF8C00",
	"#BA55D3", "#008B8B", "#B22222", "#228B22", "#8A2BE2",
	"#FF1493", "#7FFF00", "#DC143C", "#00CED1", "#FF4500",
	"#6B8E23", "#8B008B", "#4169E1", "#008080", "#8B4513",
	"#483D8B", "#00FA9A", "#8FBC8F", "#DAA520", "#7CFC00",
	"#9370DB", "#3CB371", "#FA8072", "#556B2F", "#9932CC",
	"#2E8B57", "#EE82EE", "#DE3163", "#008000",
}

type Tag struct {
	ID    int    `gorm:"primary_key;AUTO_INCREMENT"`
	Name  string `gorm:"unique"`
	Color string
}
