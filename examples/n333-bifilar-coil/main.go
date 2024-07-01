// n333-bifilar-coil creates Gerber files (and a bundled ZIP) representing
// 333 bifilar coils (https://en.wikipedia.org/wiki/Bifilar_coil).
package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime/pprof"

	_ "github.com/gmlewis/go-fonts-f/fonts/freeserif"
	. "github.com/gmlewis/go-gerber/gerber"
	"github.com/gmlewis/go-gerber/gerber/viewer"
)

var (
	step       = flag.Float64("step", 0.02, "Resolution (in radians) of the spiral")
	n          = flag.Int("n", 12, "Number of full winds in each spiral")
	gap        = flag.Float64("gap", 0.15, "Gap between traces in mm (6mil = 0.15mm)")
	trace      = flag.Float64("trace", 0.15, "Width of traces in mm")
	prefix     = flag.String("prefix", "n333-bifilar-coil", "Filename prefix for all Gerber files and zip")
	fontName   = flag.String("font", "freeserif", "Name of font to use for writing source on PCB (empty to not write)")
	view       = flag.Bool("view", false, "View the resulting design using Fyne")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

const (
	nlayers    = 334
	angleDelta = 2.0 * math.Pi / nlayers

	messageFmt = `This is a 333-layer
(on 334 layers) bifilar coil.
Trace size = %0.2fmm.
Gap size = %0.2fmm.
Each spiral has %v coils.`
)

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *n < 360 {
		flag.Usage()
		log.Fatal("N must be >= 360.")
	}

	g := New(fmt.Sprintf("%v-n%v", *prefix, *n))

	s := newSpiral()

	padD := 2.0
	startTopR, topSpiralR, endTopR := s.genSpiral(1, 0, 0)
	startTopL, topSpiralL, endTopL := s.genSpiral(1, math.Pi, 0)

	startLayerNR := map[int]Pt{}
	startLayerNL := map[int]Pt{}
	endLayerNR := map[int]Pt{}
	endLayerNL := map[int]Pt{}
	layerNSpiralR := map[int][]Pt{}
	layerNSpiralL := map[int][]Pt{}
	for n := 2; n < nlayers; n += 4 {
		af := float64((n + 2) / 4)
		startLayerNR[n], layerNSpiralR[n], endLayerNR[n] = s.genSpiral(1, af*angleDelta, 0)
		startLayerNL[n], layerNSpiralL[n], endLayerNL[n] = s.genSpiral(1, math.Pi+af*angleDelta, 0)
		startLayerNR[n+1], layerNSpiralR[n+1], endLayerNR[n+1] = s.genSpiral(-1, af*angleDelta, 0)
		startLayerNL[n+1], layerNSpiralL[n+1], endLayerNL[n+1] = s.genSpiral(-1, math.Pi+af*angleDelta, 0)

		if n+2 < nlayers {
			startLayerNR[n+2], layerNSpiralR[n+2], endLayerNR[n+2] = s.genSpiral(1, -af*angleDelta, 0)
			startLayerNL[n+2], layerNSpiralL[n+2], endLayerNL[n+2] = s.genSpiral(1, math.Pi-af*angleDelta, 0)
			startLayerNR[n+3], layerNSpiralR[n+3], endLayerNR[n+3] = s.genSpiral(-1, -af*angleDelta, 0)
			trimY := 0.0
			if n+3 == 5 { // 5L
				trimY = 0.01
			}
			startLayerNL[n+3], layerNSpiralL[n+3], endLayerNL[n+3] = s.genSpiral(-1, math.Pi-af*angleDelta, trimY)
		}
	}

	viaPadD := 0.5
	innerR := 0.5 * (*gap + viaPadD) / math.Sin(0.5*angleDelta)
	// log.Printf("innerR=%v, angleDelta=%v", innerR, angleDelta)
	// minStartAngle := (innerR + *gap + 0.5**trace + 0.5*viaPadD) * (3 * math.Pi)
	// log.Printf("innerR=%v, minStartAngle=%v/Pi=%v", innerR, minStartAngle, minStartAngle/math.Pi)
	var innerViaPts []Pt
	for i := 0; i < nlayers; i++ {
		x := innerR * math.Cos(float64(i)*angleDelta)
		y := innerR * math.Sin(float64(i)*angleDelta)
		innerViaPts = append(innerViaPts, Pt{x, y})
	}
	innerHole := map[string]int{
		"TR": 0, "TL": 167, "BR": 167, "BL": 0,
		"2R": 1, "2L": 168, "3R": 166, "3L": 333,
		"4R": 333, "4L": 166, "5R": 168, "5L": 1,
		"6R": 2, "6L": 169, "7R": 165, "7L": 332,
		"8R": 332, "8L": 165, "9R": 169, "9L": 2,
		"10R": 3, "10L": 170, "11R": 164, "11L": 331,
		"12R": 331, "12L": 164, "13R": 170, "13L": 3,
		"14R": 4, "14L": 171, "15R": 163, "15L": 330,
		"16R": 330, "16L": 163, "17R": 171, "17L": 4,
		"18R": 5, "18L": 172, "19R": 162, "19L": 329,
		"20R": 329, "20L": 162, "21R": 172, "21L": 5,
		"22R": 6, "22L": 173, "23R": 161, "23L": 328,
		"24R": 328, "24L": 161, "25R": 173, "25L": 6,
		"26R": 7, "26L": 174, "27R": 160, "27L": 327,
		"28R": 327, "28L": 160, "29R": 174, "29L": 7,
		"30R": 8, "30L": 175, "31R": 159, "31L": 326,
		"32R": 326, "32L": 159, "33R": 175, "33L": 8,
		"34R": 9, "34L": 176, "35R": 158, "35L": 325,
		"36R": 325, "36L": 158, "37R": 176, "37L": 9,
		"38R": 10, "38L": 177, "39R": 157, "39L": 324,
		"40R": 324, "40L": 157, "41R": 177, "41L": 10,
		"42R": 11, "42L": 178, "43R": 156, "43L": 323,
		"44R": 323, "44L": 156, "45R": 178, "45L": 11,
		"46R": 12, "46L": 179, "47R": 155, "47L": 322,
		"48R": 322, "48L": 155, "49R": 179, "49L": 12,
		"50R": 13, "50L": 180, "51R": 154, "51L": 321,
		"52R": 321, "52L": 154, "53R": 180, "53L": 13,
		"54R": 14, "54L": 181, "55R": 153, "55L": 320,
		"56R": 320, "56L": 153, "57R": 181, "57L": 14,
		"58R": 15, "58L": 182, "59R": 152, "59L": 319,
		"60R": 319, "60L": 152, "61R": 182, "61L": 15,
		"62R": 16, "62L": 183, "63R": 151, "63L": 318,
		"64R": 318, "64L": 151, "65R": 183, "65L": 16,
		"66R": 17, "66L": 184, "67R": 150, "67L": 317,
		"68R": 317, "68L": 150, "69R": 184, "69L": 17,
		"70R": 18, "70L": 185, "71R": 149, "71L": 316,
		"72R": 316, "72L": 149, "73R": 185, "73L": 18,
		"74R": 19, "74L": 186, "75R": 148, "75L": 315,
		"76R": 315, "76L": 148, "77R": 186, "77L": 19,
		"78R": 20, "78L": 187, "79R": 147, "79L": 314,
		"80R": 314, "80L": 147, "81R": 187, "81L": 20,
		"82R": 21, "82L": 188, "83R": 146, "83L": 313,
		"84R": 313, "84L": 146, "85R": 188, "85L": 21,
		"86R": 22, "86L": 189, "87R": 145, "87L": 312,
		"88R": 312, "88L": 145, "89R": 189, "89L": 22,
		"90R": 23, "90L": 190, "91R": 144, "91L": 311,
		"92R": 311, "92L": 144, "93R": 190, "93L": 23,
		"94R": 24, "94L": 191, "95R": 143, "95L": 310,
		"96R": 310, "96L": 143, "97R": 191, "97L": 24,
		"98R": 25, "98L": 192, "99R": 142, "99L": 309,
		"100R": 309, "100L": 142, "101R": 192, "101L": 25,
		"102R": 26, "102L": 193, "103R": 141, "103L": 308,
		"104R": 308, "104L": 141, "105R": 193, "105L": 26,
		"106R": 27, "106L": 194, "107R": 140, "107L": 307,
		"108R": 307, "108L": 140, "109R": 194, "109L": 27,
		"110R": 28, "110L": 195, "111R": 139, "111L": 306,
		"112R": 306, "112L": 139, "113R": 195, "113L": 28,
		"114R": 29, "114L": 196, "115R": 138, "115L": 305,
		"116R": 305, "116L": 138, "117R": 196, "117L": 29,
		"118R": 30, "118L": 197, "119R": 137, "119L": 304,
		"120R": 304, "120L": 137, "121R": 197, "121L": 30,
		"122R": 31, "122L": 198, "123R": 136, "123L": 303,
		"124R": 303, "124L": 136, "125R": 198, "125L": 31,
		"126R": 32, "126L": 199, "127R": 135, "127L": 302,
		"128R": 302, "128L": 135, "129R": 199, "129L": 32,
		"130R": 33, "130L": 200, "131R": 134, "131L": 301,
		"132R": 301, "132L": 134, "133R": 200, "133L": 33,
		"134R": 34, "134L": 201, "135R": 133, "135L": 300,
		"136R": 300, "136L": 133, "137R": 201, "137L": 34,
		"138R": 35, "138L": 202, "139R": 132, "139L": 299,
		"140R": 299, "140L": 132, "141R": 202, "141L": 35,
		"142R": 36, "142L": 203, "143R": 131, "143L": 298,
		"144R": 298, "144L": 131, "145R": 203, "145L": 36,
		"146R": 37, "146L": 204, "147R": 130, "147L": 297,
		"148R": 297, "148L": 130, "149R": 204, "149L": 37,
		"150R": 38, "150L": 205, "151R": 129, "151L": 296,
		"152R": 296, "152L": 129, "153R": 205, "153L": 38,
		"154R": 39, "154L": 206, "155R": 128, "155L": 295,
		"156R": 295, "156L": 128, "157R": 206, "157L": 39,
		"158R": 40, "158L": 207, "159R": 127, "159L": 294,
		"160R": 294, "160L": 127, "161R": 207, "161L": 40,
		"162R": 41, "162L": 208, "163R": 126, "163L": 293,
		"164R": 293, "164L": 126, "165R": 208, "165L": 41,
		"166R": 42, "166L": 209, "167R": 125, "167L": 292,
		"168R": 292, "168L": 125, "169R": 209, "169L": 42,
		"170R": 43, "170L": 210, "171R": 124, "171L": 291,
		"172R": 291, "172L": 124, "173R": 210, "173L": 43,
		"174R": 44, "174L": 211, "175R": 123, "175L": 290,
		"176R": 290, "176L": 123, "177R": 211, "177L": 44,
		"178R": 45, "178L": 212, "179R": 122, "179L": 289,
		"180R": 289, "180L": 122, "181R": 212, "181L": 45,
		"182R": 46, "182L": 213, "183R": 121, "183L": 288,
		"184R": 288, "184L": 121, "185R": 213, "185L": 46,
		"186R": 47, "186L": 214, "187R": 120, "187L": 287,
		"188R": 287, "188L": 120, "189R": 214, "189L": 47,
		"190R": 48, "190L": 215, "191R": 119, "191L": 286,
		"192R": 286, "192L": 119, "193R": 215, "193L": 48,
		"194R": 49, "194L": 216, "195R": 118, "195L": 285,
		"196R": 285, "196L": 118, "197R": 216, "197L": 49,
		"198R": 50, "198L": 217, "199R": 117, "199L": 284,
		"200R": 284, "200L": 117, "201R": 217, "201L": 50,
		"202R": 51, "202L": 218, "203R": 116, "203L": 283,
		"204R": 283, "204L": 116, "205R": 218, "205L": 51,
		"206R": 52, "206L": 219, "207R": 115, "207L": 282,
		"208R": 282, "208L": 115, "209R": 219, "209L": 52,
		"210R": 53, "210L": 220, "211R": 114, "211L": 281,
		"212R": 281, "212L": 114, "213R": 220, "213L": 53,
		"214R": 54, "214L": 221, "215R": 113, "215L": 280,
		"216R": 280, "216L": 113, "217R": 221, "217L": 54,
		"218R": 55, "218L": 222, "219R": 112, "219L": 279,
		"220R": 279, "220L": 112, "221R": 222, "221L": 55,
		"222R": 56, "222L": 223, "223R": 111, "223L": 278,
		"224R": 278, "224L": 111, "225R": 223, "225L": 56,
		"226R": 57, "226L": 224, "227R": 110, "227L": 277,
		"228R": 277, "228L": 110, "229R": 224, "229L": 57,
		"230R": 58, "230L": 225, "231R": 109, "231L": 276,
		"232R": 276, "232L": 109, "233R": 225, "233L": 58,
		"234R": 59, "234L": 226, "235R": 108, "235L": 275,
		"236R": 275, "236L": 108, "237R": 226, "237L": 59,
		"238R": 60, "238L": 227, "239R": 107, "239L": 274,
		"240R": 274, "240L": 107, "241R": 227, "241L": 60,
		"242R": 61, "242L": 228, "243R": 106, "243L": 273,
		"244R": 273, "244L": 106, "245R": 228, "245L": 61,
		"246R": 62, "246L": 229, "247R": 105, "247L": 272,
		"248R": 272, "248L": 105, "249R": 229, "249L": 62,
		"250R": 63, "250L": 230, "251R": 104, "251L": 271,
		"252R": 271, "252L": 104, "253R": 230, "253L": 63,
		"254R": 64, "254L": 231, "255R": 103, "255L": 270,
		"256R": 270, "256L": 103, "257R": 231, "257L": 64,
		"258R": 65, "258L": 232, "259R": 102, "259L": 269,
		"260R": 269, "260L": 102, "261R": 232, "261L": 65,
		"262R": 66, "262L": 233, "263R": 101, "263L": 268,
		"264R": 268, "264L": 101, "265R": 233, "265L": 66,
		"266R": 67, "266L": 234, "267R": 100, "267L": 267,
		"268R": 267, "268L": 100, "269R": 234, "269L": 67,
		"270R": 68, "270L": 235, "271R": 99, "271L": 266,
		"272R": 266, "272L": 99, "273R": 235, "273L": 68,
		"274R": 69, "274L": 236, "275R": 98, "275L": 265,
		"276R": 265, "276L": 98, "277R": 236, "277L": 69,
		"278R": 70, "278L": 237, "279R": 97, "279L": 264,
		"280R": 264, "280L": 97, "281R": 237, "281L": 70,
		"282R": 71, "282L": 238, "283R": 96, "283L": 263,
		"284R": 263, "284L": 96, "285R": 238, "285L": 71,
		"286R": 72, "286L": 239, "287R": 95, "287L": 262,
		"288R": 262, "288L": 95, "289R": 239, "289L": 72,
		"290R": 73, "290L": 240, "291R": 94, "291L": 261,
		"292R": 261, "292L": 94, "293R": 240, "293L": 73,
		"294R": 74, "294L": 241, "295R": 93, "295L": 260,
		"296R": 260, "296L": 93, "297R": 241, "297L": 74,
		"298R": 75, "298L": 242, "299R": 92, "299L": 259,
		"300R": 259, "300L": 92, "301R": 242, "301L": 75,
		"302R": 76, "302L": 243, "303R": 91, "303L": 258,
		"304R": 258, "304L": 91, "305R": 243, "305L": 76,
		"306R": 77, "306L": 244, "307R": 90, "307L": 257,
		"308R": 257, "308L": 90, "309R": 244, "309L": 77,
		"310R": 78, "310L": 245, "311R": 89, "311L": 256,
		"312R": 256, "312L": 89, "313R": 245, "313L": 78,
		"314R": 79, "314L": 246, "315R": 88, "315L": 255,
		"316R": 255, "316L": 88, "317R": 246, "317L": 79,
		"318R": 80, "318L": 247, "319R": 87, "319L": 254,
		"320R": 254, "320L": 87, "321R": 247, "321L": 80,
		"322R": 81, "322L": 248, "323R": 86, "323L": 253,
		"324R": 253, "324L": 86, "325R": 248, "325L": 81,
		"326R": 82, "326L": 249, "327R": 85, "327L": 252,
		"328R": 252, "328L": 85, "329R": 249, "329L": 82,
		"330R": 83, "330L": 250, "331R": 84, "331L": 251,
		"332R": 251, "332L": 84, "333R": 250, "333L": 83,
	}

	outerR := (106.0*math.Pi + float64(*n)*2.0*math.Pi + *trace + *gap) / (3.0 * math.Pi)
	outerContactPt := func(n float64) Pt {
		r := outerR + 0.5**trace + *gap + 0.5*padD
		x := r * math.Cos(n*angleDelta)
		y := r * math.Sin(n*angleDelta)
		return Pt{x, y}
	}

	var outerViaPts []Pt
	for i := 0; i < nlayers; i++ {
		pt := outerContactPt(float64(i) - 0.5)
		outerViaPts = append(outerViaPts, pt)
	}
	outerViaPts = append(outerViaPts, outerContactPt(0.0))
	outerHole := map[string]int{
		"TR": 0, "TL": 167, "BR": 166, "BL": 333,
		"2R": 1, "2L": 168, "3R": 165, "3L": 332,
		"4R": 333, "4L": 166, "5R": 167, "5L": 334,
		"6R": 2, "6L": 169, "7R": 164, "7L": 331,
		"8R": 332, "8L": 165, "9R": 168, "9L": 1,
		"10R": 3, "10L": 170, "11R": 163, "11L": 330,
		"12R": 331, "12L": 164, "13R": 169, "13L": 2,
		"14R": 4, "14L": 171, "15R": 162, "15L": 329,
		"16R": 330, "16L": 163, "17R": 170, "17L": 3,
		"18R": 5, "18L": 172, "19R": 161, "19L": 328,
		"20R": 329, "20L": 162, "21R": 171, "21L": 4,
		"22R": 6, "22L": 173, "23R": 160, "23L": 327,
		"24R": 328, "24L": 161, "25R": 172, "25L": 5,
		"26R": 7, "26L": 174, "27R": 159, "27L": 326,
		"28R": 327, "28L": 160, "29R": 173, "29L": 6,
		"30R": 8, "30L": 175, "31R": 158, "31L": 325,
		"32R": 326, "32L": 159, "33R": 174, "33L": 7,
		"34R": 9, "34L": 176, "35R": 157, "35L": 324,
		"36R": 325, "36L": 158, "37R": 175, "37L": 8,
		"38R": 10, "38L": 177, "39R": 156, "39L": 323,
		"40R": 324, "40L": 157, "41R": 176, "41L": 9,
		"42R": 11, "42L": 178, "43R": 155, "43L": 322,
		"44R": 323, "44L": 156, "45R": 177, "45L": 10,
		"46R": 12, "46L": 179, "47R": 154, "47L": 321,
		"48R": 322, "48L": 155, "49R": 178, "49L": 11,
		"50R": 13, "50L": 180, "51R": 153, "51L": 320,
		"52R": 321, "52L": 154, "53R": 179, "53L": 12,
		"54R": 14, "54L": 181, "55R": 152, "55L": 319,
		"56R": 320, "56L": 153, "57R": 180, "57L": 13,
		"58R": 15, "58L": 182, "59R": 151, "59L": 318,
		"60R": 319, "60L": 152, "61R": 181, "61L": 14,
		"62R": 16, "62L": 183, "63R": 150, "63L": 317,
		"64R": 318, "64L": 151, "65R": 182, "65L": 15,
		"66R": 17, "66L": 184, "67R": 149, "67L": 316,
		"68R": 317, "68L": 150, "69R": 183, "69L": 16,
		"70R": 18, "70L": 185, "71R": 148, "71L": 315,
		"72R": 316, "72L": 149, "73R": 184, "73L": 17,
		"74R": 19, "74L": 186, "75R": 147, "75L": 314,
		"76R": 315, "76L": 148, "77R": 185, "77L": 18,
		"78R": 20, "78L": 187, "79R": 146, "79L": 313,
		"80R": 314, "80L": 147, "81R": 186, "81L": 19,
		"82R": 21, "82L": 188, "83R": 145, "83L": 312,
		"84R": 313, "84L": 146, "85R": 187, "85L": 20,
		"86R": 22, "86L": 189, "87R": 144, "87L": 311,
		"88R": 312, "88L": 145, "89R": 188, "89L": 21,
		"90R": 23, "90L": 190, "91R": 143, "91L": 310,
		"92R": 311, "92L": 144, "93R": 189, "93L": 22,
		"94R": 24, "94L": 191, "95R": 142, "95L": 309,
		"96R": 310, "96L": 143, "97R": 190, "97L": 23,
		"98R": 25, "98L": 192, "99R": 141, "99L": 308,
		"100R": 309, "100L": 142, "101R": 191, "101L": 24,
		"102R": 26, "102L": 193, "103R": 140, "103L": 307,
		"104R": 308, "104L": 141, "105R": 192, "105L": 25,
		"106R": 27, "106L": 194, "107R": 139, "107L": 306,
		"108R": 307, "108L": 140, "109R": 193, "109L": 26,
		"110R": 28, "110L": 195, "111R": 138, "111L": 305,
		"112R": 306, "112L": 139, "113R": 194, "113L": 27,
		"114R": 29, "114L": 196, "115R": 137, "115L": 304,
		"116R": 305, "116L": 138, "117R": 195, "117L": 28,
		"118R": 30, "118L": 197, "119R": 136, "119L": 303,
		"120R": 304, "120L": 137, "121R": 196, "121L": 29,
		"122R": 31, "122L": 198, "123R": 135, "123L": 302,
		"124R": 303, "124L": 136, "125R": 197, "125L": 30,
		"126R": 32, "126L": 199, "127R": 134, "127L": 301,
		"128R": 302, "128L": 135, "129R": 198, "129L": 31,
		"130R": 33, "130L": 200, "131R": 133, "131L": 300,
		"132R": 301, "132L": 134, "133R": 199, "133L": 32,
		"134R": 34, "134L": 201, "135R": 132, "135L": 299,
		"136R": 300, "136L": 133, "137R": 200, "137L": 33,
		"138R": 35, "138L": 202, "139R": 131, "139L": 298,
		"140R": 299, "140L": 132, "141R": 201, "141L": 34,
		"142R": 36, "142L": 203, "143R": 130, "143L": 297,
		"144R": 298, "144L": 131, "145R": 202, "145L": 35,
		"146R": 37, "146L": 204, "147R": 129, "147L": 296,
		"148R": 297, "148L": 130, "149R": 203, "149L": 36,
		"150R": 38, "150L": 205, "151R": 128, "151L": 295,
		"152R": 296, "152L": 129, "153R": 204, "153L": 37,
		"154R": 39, "154L": 206, "155R": 127, "155L": 294,
		"156R": 295, "156L": 128, "157R": 205, "157L": 38,
		"158R": 40, "158L": 207, "159R": 126, "159L": 293,
		"160R": 294, "160L": 127, "161R": 206, "161L": 39,
		"162R": 41, "162L": 208, "163R": 125, "163L": 292,
		"164R": 293, "164L": 126, "165R": 207, "165L": 40,
		"166R": 42, "166L": 209, "167R": 124, "167L": 291,
		"168R": 292, "168L": 125, "169R": 208, "169L": 41,
		"170R": 43, "170L": 210, "171R": 123, "171L": 290,
		"172R": 291, "172L": 124, "173R": 209, "173L": 42,
		"174R": 44, "174L": 211, "175R": 122, "175L": 289,
		"176R": 290, "176L": 123, "177R": 210, "177L": 43,
		"178R": 45, "178L": 212, "179R": 121, "179L": 288,
		"180R": 289, "180L": 122, "181R": 211, "181L": 44,
		"182R": 46, "182L": 213, "183R": 120, "183L": 287,
		"184R": 288, "184L": 121, "185R": 212, "185L": 45,
		"186R": 47, "186L": 214, "187R": 119, "187L": 286,
		"188R": 287, "188L": 120, "189R": 213, "189L": 46,
		"190R": 48, "190L": 215, "191R": 118, "191L": 285,
		"192R": 286, "192L": 119, "193R": 214, "193L": 47,
		"194R": 49, "194L": 216, "195R": 117, "195L": 284,
		"196R": 285, "196L": 118, "197R": 215, "197L": 48,
		"198R": 50, "198L": 217, "199R": 116, "199L": 283,
		"200R": 284, "200L": 117, "201R": 216, "201L": 49,
		"202R": 51, "202L": 218, "203R": 115, "203L": 282,
		"204R": 283, "204L": 116, "205R": 217, "205L": 50,
		"206R": 52, "206L": 219, "207R": 114, "207L": 281,
		"208R": 282, "208L": 115, "209R": 218, "209L": 51,
		"210R": 53, "210L": 220, "211R": 113, "211L": 280,
		"212R": 281, "212L": 114, "213R": 219, "213L": 52,
		"214R": 54, "214L": 221, "215R": 112, "215L": 279,
		"216R": 280, "216L": 113, "217R": 220, "217L": 53,
		"218R": 55, "218L": 222, "219R": 111, "219L": 278,
		"220R": 279, "220L": 112, "221R": 221, "221L": 54,
		"222R": 56, "222L": 223, "223R": 110, "223L": 277,
		"224R": 278, "224L": 111, "225R": 222, "225L": 55,
		"226R": 57, "226L": 224, "227R": 109, "227L": 276,
		"228R": 277, "228L": 110, "229R": 223, "229L": 56,
		"230R": 58, "230L": 225, "231R": 108, "231L": 275,
		"232R": 276, "232L": 109, "233R": 224, "233L": 57,
		"234R": 59, "234L": 226, "235R": 107, "235L": 274,
		"236R": 275, "236L": 108, "237R": 225, "237L": 58,
		"238R": 60, "238L": 227, "239R": 106, "239L": 273,
		"240R": 274, "240L": 107, "241R": 226, "241L": 59,
		"242R": 61, "242L": 228, "243R": 105, "243L": 272,
		"244R": 273, "244L": 106, "245R": 227, "245L": 60,
		"246R": 62, "246L": 229, "247R": 104, "247L": 271,
		"248R": 272, "248L": 105, "249R": 228, "249L": 61,
		"250R": 63, "250L": 230, "251R": 103, "251L": 270,
		"252R": 271, "252L": 104, "253R": 229, "253L": 62,
		"254R": 64, "254L": 231, "255R": 102, "255L": 269,
		"256R": 270, "256L": 103, "257R": 230, "257L": 63,
		"258R": 65, "258L": 232, "259R": 101, "259L": 268,
		"260R": 269, "260L": 102, "261R": 231, "261L": 64,
		"262R": 66, "262L": 233, "263R": 100, "263L": 267,
		"264R": 268, "264L": 101, "265R": 232, "265L": 65,
		"266R": 67, "266L": 234, "267R": 99, "267L": 266,
		"268R": 267, "268L": 100, "269R": 233, "269L": 66,
		"270R": 68, "270L": 235, "271R": 98, "271L": 265,
		"272R": 266, "272L": 99, "273R": 234, "273L": 67,
		"274R": 69, "274L": 236, "275R": 97, "275L": 264,
		"276R": 265, "276L": 98, "277R": 235, "277L": 68,
		"278R": 70, "278L": 237, "279R": 96, "279L": 263,
		"280R": 264, "280L": 97, "281R": 236, "281L": 69,
		"282R": 71, "282L": 238, "283R": 95, "283L": 262,
		"284R": 263, "284L": 96, "285R": 237, "285L": 70,
		"286R": 72, "286L": 239, "287R": 94, "287L": 261,
		"288R": 262, "288L": 95, "289R": 238, "289L": 71,
		"290R": 73, "290L": 240, "291R": 93, "291L": 260,
		"292R": 261, "292L": 94, "293R": 239, "293L": 72,
		"294R": 74, "294L": 241, "295R": 92, "295L": 259,
		"296R": 260, "296L": 93, "297R": 240, "297L": 73,
		"298R": 75, "298L": 242, "299R": 91, "299L": 258,
		"300R": 259, "300L": 92, "301R": 241, "301L": 74,
		"302R": 76, "302L": 243, "303R": 90, "303L": 257,
		"304R": 258, "304L": 91, "305R": 242, "305L": 75,
		"306R": 77, "306L": 244, "307R": 89, "307L": 256,
		"308R": 257, "308L": 90, "309R": 243, "309L": 76,
		"310R": 78, "310L": 245, "311R": 88, "311L": 255,
		"312R": 256, "312L": 89, "313R": 244, "313L": 77,
		"314R": 79, "314L": 246, "315R": 87, "315L": 254,
		"316R": 255, "316L": 88, "317R": 245, "317L": 78,
		"318R": 80, "318L": 247, "319R": 86, "319L": 253,
		"320R": 254, "320L": 87, "321R": 246, "321L": 79,
		"322R": 81, "322L": 248, "323R": 85, "323L": 252,
		"324R": 253, "324L": 86, "325R": 247, "325L": 80,
		"326R": 82, "326L": 249, "327R": 84, "327L": 251,
		"328R": 252, "328L": 85, "329R": 248, "329L": 81,
		"330R": 83, "330L": 250, "331R": 83, "331L": 250,
		"332R": 251, "332L": 84, "333R": 249, "333L": 82,
	}

	drill := g.Drill()
	for _, pt := range innerViaPts {
		const viaDrillD = 0.25
		drill.Add(Circle(pt, viaDrillD))
	}
	for _, pt := range outerViaPts {
		const drillD = 1.0
		drill.Add(Circle(pt, drillD))
	}

	padLine := func(pt1, pt2 Pt) *LineT {
		return Line(pt1[0], pt1[1], pt2[0], pt2[1], CircleShape, *trace)
	}
	addVias := func(layer *Layer) {
		for _, pt := range innerViaPts {
			layer.Add(Circle(pt, viaPadD))
		}
		for _, pt := range outerViaPts {
			layer.Add(Circle(pt, padD))
		}
	}

	top := g.TopCopper()
	top.Add(
		Polygon(Pt{0, 0}, true, topSpiralR, 0.0),
		Polygon(Pt{0, 0}, true, topSpiralL, 0.0),
		padLine(startTopR, innerViaPts[innerHole["TR"]]),
		padLine(startTopL, innerViaPts[innerHole["TL"]]),
		padLine(endTopR, outerViaPts[outerHole["TR"]]),
		padLine(endTopL, outerViaPts[outerHole["TL"]]),
	)
	addVias(top)

	topMask := g.TopSolderMask()
	addVias(topMask)

	for n := 2; n < nlayers; n++ {
		if n > 21 {
			continue // ONLY FOR DIAGRAM
		}
		nr := fmt.Sprintf("%vR", n)
		nl := fmt.Sprintf("%vL", n)
		layer := g.LayerN(n)
		layer.Add(
			Polygon(Pt{0, 0}, true, layerNSpiralR[n], 0.0),
			Polygon(Pt{0, 0}, true, layerNSpiralL[n], 0.0),
			padLine(startLayerNR[n], innerViaPts[innerHole[nr]]),
			padLine(startLayerNL[n], innerViaPts[innerHole[nl]]),
			padLine(endLayerNR[n], outerViaPts[outerHole[nr]]),
			padLine(endLayerNL[n], outerViaPts[outerHole[nl]]),
		)
		addVias(layer)
	}

	bottom := g.BottomCopper()
	bottom.Add(
		padLine(innerViaPts[innerHole["TR"]], outerViaPts[outerHole["4R"]]),
		padLine(innerViaPts[innerHole["TL"]], outerViaPts[outerHole["4L"]]),
	)
	addVias(bottom)

	bottomMask := g.BottomSolderMask()
	addVias(bottomMask)

	outline := g.Outline()
	r := 0.5*s.size + padD + *trace
	outline.Add(
		Arc(Pt{0, 0}, r, CircleShape, 1, 1, 0, 360, 0.1),
	)
	fmt.Printf("n=%v: (%.2f,%.2f)\n", *n, 2*r, 2*r)

	if *fontName != "" {
		pts := 15.9
		labelSize := 0.45
		outerLabelSize := 3.0
		message := fmt.Sprintf(messageFmt, *trace, *gap, *n)

		innerLabel := func(label string) *TextT {
			num := float64(innerHole[label])
			r := innerR - viaPadD
			x := r * math.Cos(num*angleDelta)
			y := r * math.Sin(num*angleDelta)
			return Text(x, y, 1.0, label, *fontName, labelSize, &Center)
		}
		innerLabel2 := func(label string) *TextT {
			num := float64(innerHole[label])
			r := innerR + viaPadD
			x := r * math.Cos(num*angleDelta)
			y := r * math.Sin(num*angleDelta)
			return Text(x, y, 1.0, label, *fontName, labelSize, &Center)
		}
		outerLabel := func(label string) *TextT {
			num := float64(outerHole[label])
			num -= 0.5
			r := outerR + 0.5**trace + *gap + 0.5*padD
			x := r * math.Cos((0.3+num)*angleDelta)
			y := r * math.Sin((0.3+num)*angleDelta)
			return Text(x, y, 1.0, label, *fontName, outerLabelSize, &Center)
		}
		outerLabel2 := func(label string) *TextT {
			num := float64(outerHole[label])
			num -= 0.5
			r := outerR + 0.5**trace + *gap + 0.5*padD
			x := r * math.Cos((-0.3+num)*angleDelta)
			y := r * math.Sin((-0.3+num)*angleDelta)
			return Text(x, y, 1.0, label, *fontName, outerLabelSize, &Center)
		}
		outerLabel3 := func(label string, num float64) *TextT {
			r := outerR + 0.5**trace + *gap
			x := r * math.Cos(num*angleDelta)
			y := r * math.Sin(num*angleDelta)
			return Text(x, y, 1.0, label, *fontName, outerLabelSize, &CenterRight)
		}

		tss := g.TopSilkscreen()
		tss.Add(
			Text(0, 0, 1.0, message, *fontName, pts, &Center),
			innerLabel("TR"), innerLabel("TL"), innerLabel2("BR"), innerLabel2("BL"),
			innerLabel("2R"), innerLabel("2L"), innerLabel("3R"), innerLabel("3L"),
			innerLabel2("4R"), innerLabel2("4L"), innerLabel2("5R"), innerLabel2("5L"),
			innerLabel("6R"), innerLabel("6L"), innerLabel("7R"), innerLabel("7L"),
			innerLabel2("8R"), innerLabel2("8L"), innerLabel2("9R"), innerLabel2("9L"),
			innerLabel("10R"), innerLabel("10L"), innerLabel("11R"), innerLabel("11L"),
			innerLabel2("12R"), innerLabel2("12L"), innerLabel2("13R"), innerLabel2("13L"),
			innerLabel("14R"), innerLabel("14L"), innerLabel("15R"), innerLabel("15L"),
			innerLabel2("16R"), innerLabel2("16L"), innerLabel2("17R"), innerLabel2("17L"),
			innerLabel("18R"), innerLabel("18L"), innerLabel("19R"), innerLabel("19L"),
			innerLabel2("20R"), innerLabel2("20L"), innerLabel2("21R"), innerLabel2("21L"),
			innerLabel("22R"), innerLabel("22L"), innerLabel("23R"), innerLabel("23L"),
			innerLabel2("24R"), innerLabel2("24L"), innerLabel2("25R"), innerLabel2("25L"),
			innerLabel("26R"), innerLabel("26L"), innerLabel("27R"), innerLabel("27L"),
			innerLabel2("28R"), innerLabel2("28L"), innerLabel2("29R"), innerLabel2("29L"),
			innerLabel("30R"), innerLabel("30L"), innerLabel("31R"), innerLabel("31L"),
			innerLabel2("32R"), innerLabel2("32L"), innerLabel2("33R"), innerLabel2("33L"),
			innerLabel("34R"), innerLabel("34L"), innerLabel("35R"), innerLabel("35L"),
			innerLabel2("36R"), innerLabel2("36L"), innerLabel2("37R"), innerLabel2("37L"),
			innerLabel("38R"), innerLabel("38L"), innerLabel("39R"), innerLabel("39L"),
			innerLabel2("40R"), innerLabel2("40L"), innerLabel2("41R"), innerLabel2("41L"),
			innerLabel("42R"), innerLabel("42L"), innerLabel("43R"), innerLabel("43L"),
			innerLabel2("44R"), innerLabel2("44L"), innerLabel2("45R"), innerLabel2("45L"),
			innerLabel("46R"), innerLabel("46L"), innerLabel("47R"), innerLabel("47L"),
			innerLabel2("48R"), innerLabel2("48L"), innerLabel2("49R"), innerLabel2("49L"),
			innerLabel("50R"), innerLabel("50L"), innerLabel("51R"), innerLabel("51L"),
			innerLabel2("52R"), innerLabel2("52L"), innerLabel2("53R"), innerLabel2("53L"),
			innerLabel("54R"), innerLabel("54L"), innerLabel("55R"), innerLabel("55L"),
			innerLabel2("56R"), innerLabel2("56L"), innerLabel2("57R"), innerLabel2("57L"),
			innerLabel("58R"), innerLabel("58L"), innerLabel("59R"), innerLabel("59L"),
			innerLabel2("60R"), innerLabel2("60L"), innerLabel2("61R"), innerLabel2("61L"),
			innerLabel("62R"), innerLabel("62L"), innerLabel("63R"), innerLabel("63L"),
			innerLabel2("64R"), innerLabel2("64L"), innerLabel2("65R"), innerLabel2("65L"),
			innerLabel("66R"), innerLabel("66L"), innerLabel("67R"), innerLabel("67L"),
			innerLabel2("68R"), innerLabel2("68L"), innerLabel2("69R"), innerLabel2("69L"),
			innerLabel("70R"), innerLabel("70L"), innerLabel("71R"), innerLabel("71L"),
			innerLabel2("72R"), innerLabel2("72L"), innerLabel2("73R"), innerLabel2("73L"),
			innerLabel("74R"), innerLabel("74L"), innerLabel("75R"), innerLabel("75L"),
			innerLabel2("76R"), innerLabel2("76L"), innerLabel2("77R"), innerLabel2("77L"),
			innerLabel("78R"), innerLabel("78L"), innerLabel("79R"), innerLabel("79L"),
			innerLabel2("80R"), innerLabel2("80L"), innerLabel2("81R"), innerLabel2("81L"),
			innerLabel("82R"), innerLabel("82L"), innerLabel("83R"), innerLabel("83L"),
			innerLabel2("84R"), innerLabel2("84L"), innerLabel2("85R"), innerLabel2("85L"),
			innerLabel("86R"), innerLabel("86L"), innerLabel("87R"), innerLabel("87L"),
			innerLabel2("88R"), innerLabel2("88L"), innerLabel2("89R"), innerLabel2("89L"),
			innerLabel("90R"), innerLabel("90L"), innerLabel("91R"), innerLabel("91L"),
			innerLabel2("92R"), innerLabel2("92L"), innerLabel2("93R"), innerLabel2("93L"),
			innerLabel("94R"), innerLabel("94L"), innerLabel("95R"), innerLabel("95L"),
			innerLabel2("96R"), innerLabel2("96L"), innerLabel2("97R"), innerLabel2("97L"),
			innerLabel("98R"), innerLabel("98L"), innerLabel("99R"), innerLabel("99L"),
			innerLabel2("100R"), innerLabel2("100L"), innerLabel2("101R"), innerLabel2("101L"),
			innerLabel("102R"), innerLabel("102L"), innerLabel("103R"), innerLabel("103L"),
			innerLabel2("104R"), innerLabel2("104L"), innerLabel2("105R"), innerLabel2("105L"),
			innerLabel("106R"), innerLabel("106L"), innerLabel("107R"), innerLabel("107L"),
			innerLabel2("108R"), innerLabel2("108L"), innerLabel2("109R"), innerLabel2("109L"),
			innerLabel("110R"), innerLabel("110L"), innerLabel("111R"), innerLabel("111L"),
			innerLabel2("112R"), innerLabel2("112L"), innerLabel2("113R"), innerLabel2("113L"),
			innerLabel("114R"), innerLabel("114L"), innerLabel("115R"), innerLabel("115L"),
			innerLabel2("116R"), innerLabel2("116L"), innerLabel2("117R"), innerLabel2("117L"),
			innerLabel("118R"), innerLabel("118L"), innerLabel("119R"), innerLabel("119L"),
			innerLabel2("120R"), innerLabel2("120L"), innerLabel2("121R"), innerLabel2("121L"),
			innerLabel("122R"), innerLabel("122L"), innerLabel("123R"), innerLabel("123L"),
			innerLabel2("124R"), innerLabel2("124L"), innerLabel2("125R"), innerLabel2("125L"),
			innerLabel("126R"), innerLabel("126L"), innerLabel("127R"), innerLabel("127L"),
			innerLabel2("128R"), innerLabel2("128L"), innerLabel2("129R"), innerLabel2("129L"),
			innerLabel("130R"), innerLabel("130L"), innerLabel("131R"), innerLabel("131L"),
			innerLabel2("132R"), innerLabel2("132L"), innerLabel2("133R"), innerLabel2("133L"),
			innerLabel("134R"), innerLabel("134L"), innerLabel("135R"), innerLabel("135L"),
			innerLabel2("136R"), innerLabel2("136L"), innerLabel2("137R"), innerLabel2("137L"),
			innerLabel("138R"), innerLabel("138L"), innerLabel("139R"), innerLabel("139L"),
			innerLabel2("140R"), innerLabel2("140L"), innerLabel2("141R"), innerLabel2("141L"),
			innerLabel("142R"), innerLabel("142L"), innerLabel("143R"), innerLabel("143L"),
			innerLabel2("144R"), innerLabel2("144L"), innerLabel2("145R"), innerLabel2("145L"),
			innerLabel("146R"), innerLabel("146L"), innerLabel("147R"), innerLabel("147L"),
			innerLabel2("148R"), innerLabel2("148L"), innerLabel2("149R"), innerLabel2("149L"),
			innerLabel("150R"), innerLabel("150L"), innerLabel("151R"), innerLabel("151L"),
			innerLabel2("152R"), innerLabel2("152L"), innerLabel2("153R"), innerLabel2("153L"),
			innerLabel("154R"), innerLabel("154L"), innerLabel("155R"), innerLabel("155L"),
			innerLabel2("156R"), innerLabel2("156L"), innerLabel2("157R"), innerLabel2("157L"),
			innerLabel("158R"), innerLabel("158L"), innerLabel("159R"), innerLabel("159L"),
			innerLabel2("160R"), innerLabel2("160L"), innerLabel2("161R"), innerLabel2("161L"),
			innerLabel("162R"), innerLabel("162L"), innerLabel("163R"), innerLabel("163L"),
			innerLabel2("164R"), innerLabel2("164L"), innerLabel2("165R"), innerLabel2("165L"),
			innerLabel("166R"), innerLabel("166L"), innerLabel("167R"), innerLabel("167L"),
			innerLabel2("168R"), innerLabel2("168L"), innerLabel2("169R"), innerLabel2("169L"),
			innerLabel("170R"), innerLabel("170L"), innerLabel("171R"), innerLabel("171L"),
			innerLabel2("172R"), innerLabel2("172L"), innerLabel2("173R"), innerLabel2("173L"),
			innerLabel("174R"), innerLabel("174L"), innerLabel("175R"), innerLabel("175L"),
			innerLabel2("176R"), innerLabel2("176L"), innerLabel2("177R"), innerLabel2("177L"),
			innerLabel("178R"), innerLabel("178L"), innerLabel("179R"), innerLabel("179L"),
			innerLabel2("180R"), innerLabel2("180L"), innerLabel2("181R"), innerLabel2("181L"),
			innerLabel("182R"), innerLabel("182L"), innerLabel("183R"), innerLabel("183L"),
			innerLabel2("184R"), innerLabel2("184L"), innerLabel2("185R"), innerLabel2("185L"),
			innerLabel("186R"), innerLabel("186L"), innerLabel("187R"), innerLabel("187L"),
			innerLabel2("188R"), innerLabel2("188L"), innerLabel2("189R"), innerLabel2("189L"),
			innerLabel("190R"), innerLabel("190L"), innerLabel("191R"), innerLabel("191L"),
			innerLabel2("192R"), innerLabel2("192L"), innerLabel2("193R"), innerLabel2("193L"),
			innerLabel("194R"), innerLabel("194L"), innerLabel("195R"), innerLabel("195L"),
			innerLabel2("196R"), innerLabel2("196L"), innerLabel2("197R"), innerLabel2("197L"),
			innerLabel("198R"), innerLabel("198L"), innerLabel("199R"), innerLabel("199L"),
			innerLabel2("200R"), innerLabel2("200L"), innerLabel2("201R"), innerLabel2("201L"),
			innerLabel("202R"), innerLabel("202L"), innerLabel("203R"), innerLabel("203L"),
			innerLabel2("204R"), innerLabel2("204L"), innerLabel2("205R"), innerLabel2("205L"),
			innerLabel("206R"), innerLabel("206L"), innerLabel("207R"), innerLabel("207L"),
			innerLabel2("208R"), innerLabel2("208L"), innerLabel2("209R"), innerLabel2("209L"),
			innerLabel("210R"), innerLabel("210L"), innerLabel("211R"), innerLabel("211L"),
			innerLabel2("212R"), innerLabel2("212L"), innerLabel2("213R"), innerLabel2("213L"),
			innerLabel("214R"), innerLabel("214L"), innerLabel("215R"), innerLabel("215L"),
			innerLabel2("216R"), innerLabel2("216L"), innerLabel2("217R"), innerLabel2("217L"),
			innerLabel("218R"), innerLabel("218L"), innerLabel("219R"), innerLabel("219L"),
			innerLabel2("220R"), innerLabel2("220L"), innerLabel2("221R"), innerLabel2("221L"),
			innerLabel("222R"), innerLabel("222L"), innerLabel("223R"), innerLabel("223L"),
			innerLabel2("224R"), innerLabel2("224L"), innerLabel2("225R"), innerLabel2("225L"),
			innerLabel("226R"), innerLabel("226L"), innerLabel("227R"), innerLabel("227L"),
			innerLabel2("228R"), innerLabel2("228L"), innerLabel2("229R"), innerLabel2("229L"),
			innerLabel("230R"), innerLabel("230L"), innerLabel("231R"), innerLabel("231L"),
			innerLabel2("232R"), innerLabel2("232L"), innerLabel2("233R"), innerLabel2("233L"),
			innerLabel("234R"), innerLabel("234L"), innerLabel("235R"), innerLabel("235L"),
			innerLabel2("236R"), innerLabel2("236L"), innerLabel2("237R"), innerLabel2("237L"),
			innerLabel("238R"), innerLabel("238L"), innerLabel("239R"), innerLabel("239L"),
			innerLabel2("240R"), innerLabel2("240L"), innerLabel2("241R"), innerLabel2("241L"),
			innerLabel("242R"), innerLabel("242L"), innerLabel("243R"), innerLabel("243L"),
			innerLabel2("244R"), innerLabel2("244L"), innerLabel2("245R"), innerLabel2("245L"),
			innerLabel("246R"), innerLabel("246L"), innerLabel("247R"), innerLabel("247L"),
			innerLabel2("248R"), innerLabel2("248L"), innerLabel2("249R"), innerLabel2("249L"),
			innerLabel("250R"), innerLabel("250L"), innerLabel("251R"), innerLabel("251L"),
			innerLabel2("252R"), innerLabel2("252L"), innerLabel2("253R"), innerLabel2("253L"),
			innerLabel("254R"), innerLabel("254L"), innerLabel("255R"), innerLabel("255L"),
			innerLabel2("256R"), innerLabel2("256L"), innerLabel2("257R"), innerLabel2("257L"),
			innerLabel("258R"), innerLabel("258L"), innerLabel("259R"), innerLabel("259L"),
			innerLabel2("260R"), innerLabel2("260L"), innerLabel2("261R"), innerLabel2("261L"),
			innerLabel("262R"), innerLabel("262L"), innerLabel("263R"), innerLabel("263L"),
			innerLabel2("264R"), innerLabel2("264L"), innerLabel2("265R"), innerLabel2("265L"),
			innerLabel("266R"), innerLabel("266L"), innerLabel("267R"), innerLabel("267L"),
			innerLabel2("268R"), innerLabel2("268L"), innerLabel2("269R"), innerLabel2("269L"),
			innerLabel("270R"), innerLabel("270L"), innerLabel("271R"), innerLabel("271L"),
			innerLabel2("272R"), innerLabel2("272L"), innerLabel2("273R"), innerLabel2("273L"),
			innerLabel("274R"), innerLabel("274L"), innerLabel("275R"), innerLabel("275L"),
			innerLabel2("276R"), innerLabel2("276L"), innerLabel2("277R"), innerLabel2("277L"),
			innerLabel("278R"), innerLabel("278L"), innerLabel("279R"), innerLabel("279L"),
			innerLabel2("280R"), innerLabel2("280L"), innerLabel2("281R"), innerLabel2("281L"),
			innerLabel("282R"), innerLabel("282L"), innerLabel("283R"), innerLabel("283L"),
			innerLabel2("284R"), innerLabel2("284L"), innerLabel2("285R"), innerLabel2("285L"),
			innerLabel("286R"), innerLabel("286L"), innerLabel("287R"), innerLabel("287L"),
			innerLabel2("288R"), innerLabel2("288L"), innerLabel2("289R"), innerLabel2("289L"),
			innerLabel("290R"), innerLabel("290L"), innerLabel("291R"), innerLabel("291L"),
			innerLabel2("292R"), innerLabel2("292L"), innerLabel2("293R"), innerLabel2("293L"),
			innerLabel("294R"), innerLabel("294L"), innerLabel("295R"), innerLabel("295L"),
			innerLabel2("296R"), innerLabel2("296L"), innerLabel2("297R"), innerLabel2("297L"),
			innerLabel("298R"), innerLabel("298L"), innerLabel("299R"), innerLabel("299L"),
			innerLabel2("300R"), innerLabel2("300L"), innerLabel2("301R"), innerLabel2("301L"),
			innerLabel("302R"), innerLabel("302L"), innerLabel("303R"), innerLabel("303L"),
			innerLabel2("304R"), innerLabel2("304L"), innerLabel2("305R"), innerLabel2("305L"),
			innerLabel("306R"), innerLabel("306L"), innerLabel("307R"), innerLabel("307L"),
			innerLabel2("308R"), innerLabel2("308L"), innerLabel2("309R"), innerLabel2("309L"),
			innerLabel("310R"), innerLabel("310L"), innerLabel("311R"), innerLabel("311L"),
			innerLabel2("312R"), innerLabel2("312L"), innerLabel2("313R"), innerLabel2("313L"),
			innerLabel("314R"), innerLabel("314L"), innerLabel("315R"), innerLabel("315L"),
			innerLabel2("316R"), innerLabel2("316L"), innerLabel2("317R"), innerLabel2("317L"),
			innerLabel("318R"), innerLabel("318L"), innerLabel("319R"), innerLabel("319L"),
			innerLabel2("320R"), innerLabel2("320L"), innerLabel2("321R"), innerLabel2("321L"),
			innerLabel("322R"), innerLabel("322L"), innerLabel("323R"), innerLabel("323L"),
			innerLabel2("324R"), innerLabel2("324L"), innerLabel2("325R"), innerLabel2("325L"),
			innerLabel("326R"), innerLabel("326L"), innerLabel("327R"), innerLabel("327L"),
			innerLabel2("328R"), innerLabel2("328L"), innerLabel2("329R"), innerLabel2("329L"),
			innerLabel("330R"), innerLabel("330L"), innerLabel("331R"), innerLabel("331L"),
			innerLabel2("332R"), innerLabel2("332L"), innerLabel2("333R"), innerLabel2("333L"),

			outerLabel3("TR", -0.5), outerLabel("TL"), outerLabel("BR"), outerLabel("BL"),
			outerLabel("2R"), outerLabel("2L"), outerLabel("3R"), outerLabel("3L"),
			outerLabel2("4R"), outerLabel2("4L"), outerLabel2("5R"), outerLabel3("5L", 0.0),
			outerLabel("6R"), outerLabel("6L"), outerLabel("7R"), outerLabel("7L"),
			outerLabel2("8R"), outerLabel2("8L"), outerLabel2("9R"), outerLabel2("9L"),
			outerLabel("10R"), outerLabel("10L"), outerLabel("11R"), outerLabel("11L"),
			outerLabel2("12R"), outerLabel2("12L"), outerLabel2("13R"), outerLabel2("13L"),
			outerLabel("14R"), outerLabel("14L"), outerLabel("15R"), outerLabel("15L"),
			outerLabel2("16R"), outerLabel2("16L"), outerLabel2("17R"), outerLabel2("17L"),
			outerLabel("18R"), outerLabel("18L"), outerLabel("19R"), outerLabel("19L"),
			outerLabel2("20R"), outerLabel2("20L"), outerLabel2("21R"), outerLabel2("21L"),
			outerLabel("22R"), outerLabel("22L"), outerLabel("23R"), outerLabel("23L"),
			outerLabel2("24R"), outerLabel2("24L"), outerLabel2("25R"), outerLabel2("25L"),
			outerLabel("26R"), outerLabel("26L"), outerLabel("27R"), outerLabel("27L"),
			outerLabel2("28R"), outerLabel2("28L"), outerLabel2("29R"), outerLabel2("29L"),
			outerLabel("30R"), outerLabel("30L"), outerLabel("31R"), outerLabel("31L"),
			outerLabel2("32R"), outerLabel2("32L"), outerLabel2("33R"), outerLabel2("33L"),
			outerLabel("34R"), outerLabel("34L"), outerLabel("35R"), outerLabel("35L"),
			outerLabel2("36R"), outerLabel2("36L"), outerLabel2("37R"), outerLabel2("37L"),
			outerLabel("38R"), outerLabel("38L"), outerLabel("39R"), outerLabel("39L"),
			outerLabel2("40R"), outerLabel2("40L"), outerLabel2("41R"), outerLabel2("41L"),
			outerLabel("42R"), outerLabel("42L"), outerLabel("43R"), outerLabel("43L"),
			outerLabel2("44R"), outerLabel2("44L"), outerLabel2("45R"), outerLabel2("45L"),
			outerLabel("46R"), outerLabel("46L"), outerLabel("47R"), outerLabel("47L"),
			outerLabel2("48R"), outerLabel2("48L"), outerLabel2("49R"), outerLabel2("49L"),
			outerLabel("50R"), outerLabel("50L"), outerLabel("51R"), outerLabel("51L"),
			outerLabel2("52R"), outerLabel2("52L"), outerLabel2("53R"), outerLabel2("53L"),
			outerLabel("54R"), outerLabel("54L"), outerLabel("55R"), outerLabel("55L"),
			outerLabel2("56R"), outerLabel2("56L"), outerLabel2("57R"), outerLabel2("57L"),
			outerLabel("58R"), outerLabel("58L"), outerLabel("59R"), outerLabel("59L"),
			outerLabel2("60R"), outerLabel2("60L"), outerLabel2("61R"), outerLabel2("61L"),
			outerLabel("62R"), outerLabel("62L"), outerLabel("63R"), outerLabel("63L"),
			outerLabel2("64R"), outerLabel2("64L"), outerLabel2("65R"), outerLabel2("65L"),
			outerLabel("66R"), outerLabel("66L"), outerLabel("67R"), outerLabel("67L"),
			outerLabel2("68R"), outerLabel2("68L"), outerLabel2("69R"), outerLabel2("69L"),
			outerLabel("70R"), outerLabel("70L"), outerLabel("71R"), outerLabel("71L"),
			outerLabel2("72R"), outerLabel2("72L"), outerLabel2("73R"), outerLabel2("73L"),
			outerLabel("74R"), outerLabel("74L"), outerLabel("75R"), outerLabel("75L"),
			outerLabel2("76R"), outerLabel2("76L"), outerLabel2("77R"), outerLabel2("77L"),
			outerLabel("78R"), outerLabel("78L"), outerLabel("79R"), outerLabel("79L"),
			outerLabel2("80R"), outerLabel2("80L"), outerLabel2("81R"), outerLabel2("81L"),
			outerLabel("82R"), outerLabel("82L"), outerLabel("83R"), outerLabel("83L"),
			outerLabel2("84R"), outerLabel2("84L"), outerLabel2("85R"), outerLabel2("85L"),
			outerLabel("86R"), outerLabel("86L"), outerLabel("87R"), outerLabel("87L"),
			outerLabel2("88R"), outerLabel2("88L"), outerLabel2("89R"), outerLabel2("89L"),
			outerLabel("90R"), outerLabel("90L"), outerLabel("91R"), outerLabel("91L"),
			outerLabel2("92R"), outerLabel2("92L"), outerLabel2("93R"), outerLabel2("93L"),
			outerLabel("94R"), outerLabel("94L"), outerLabel("95R"), outerLabel("95L"),
			outerLabel2("96R"), outerLabel2("96L"), outerLabel2("97R"), outerLabel2("97L"),
			outerLabel("98R"), outerLabel("98L"), outerLabel("99R"), outerLabel("99L"),
			outerLabel2("100R"), outerLabel2("100L"), outerLabel2("101R"), outerLabel2("101L"),
			outerLabel("102R"), outerLabel("102L"), outerLabel("103R"), outerLabel("103L"),
			outerLabel2("104R"), outerLabel2("104L"), outerLabel2("105R"), outerLabel2("105L"),
			outerLabel("106R"), outerLabel("106L"), outerLabel("107R"), outerLabel("107L"),
			outerLabel2("108R"), outerLabel2("108L"), outerLabel2("109R"), outerLabel2("109L"),
			outerLabel("110R"), outerLabel("110L"), outerLabel("111R"), outerLabel("111L"),
			outerLabel2("112R"), outerLabel2("112L"), outerLabel2("113R"), outerLabel2("113L"),
			outerLabel("114R"), outerLabel("114L"), outerLabel("115R"), outerLabel("115L"),
			outerLabel2("116R"), outerLabel2("116L"), outerLabel2("117R"), outerLabel2("117L"),
			outerLabel("118R"), outerLabel("118L"), outerLabel("119R"), outerLabel("119L"),
			outerLabel2("120R"), outerLabel2("120L"), outerLabel2("121R"), outerLabel2("121L"),
			outerLabel("122R"), outerLabel("122L"), outerLabel("123R"), outerLabel("123L"),
			outerLabel2("124R"), outerLabel2("124L"), outerLabel2("125R"), outerLabel2("125L"),
			outerLabel("126R"), outerLabel("126L"), outerLabel("127R"), outerLabel("127L"),
			outerLabel2("128R"), outerLabel2("128L"), outerLabel2("129R"), outerLabel2("129L"),
			outerLabel("130R"), outerLabel("130L"), outerLabel("131R"), outerLabel("131L"),
			outerLabel2("132R"), outerLabel2("132L"), outerLabel2("133R"), outerLabel2("133L"),
			outerLabel("134R"), outerLabel("134L"), outerLabel("135R"), outerLabel("135L"),
			outerLabel2("136R"), outerLabel2("136L"), outerLabel2("137R"), outerLabel2("137L"),
			outerLabel("138R"), outerLabel("138L"), outerLabel("139R"), outerLabel("139L"),
			outerLabel2("140R"), outerLabel2("140L"), outerLabel2("141R"), outerLabel2("141L"),
			outerLabel("142R"), outerLabel("142L"), outerLabel("143R"), outerLabel("143L"),
			outerLabel2("144R"), outerLabel2("144L"), outerLabel2("145R"), outerLabel2("145L"),
			outerLabel("146R"), outerLabel("146L"), outerLabel("147R"), outerLabel("147L"),
			outerLabel2("148R"), outerLabel2("148L"), outerLabel2("149R"), outerLabel2("149L"),
			outerLabel("150R"), outerLabel("150L"), outerLabel("151R"), outerLabel("151L"),
			outerLabel2("152R"), outerLabel2("152L"), outerLabel2("153R"), outerLabel2("153L"),
			outerLabel("154R"), outerLabel("154L"), outerLabel("155R"), outerLabel("155L"),
			outerLabel2("156R"), outerLabel2("156L"), outerLabel2("157R"), outerLabel2("157L"),
			outerLabel("158R"), outerLabel("158L"), outerLabel("159R"), outerLabel("159L"),
			outerLabel2("160R"), outerLabel2("160L"), outerLabel2("161R"), outerLabel2("161L"),
			outerLabel("162R"), outerLabel("162L"), outerLabel("163R"), outerLabel("163L"),
			outerLabel2("164R"), outerLabel2("164L"), outerLabel2("165R"), outerLabel2("165L"),
			outerLabel("166R"), outerLabel("166L"), outerLabel("167R"), outerLabel("167L"),
			outerLabel2("168R"), outerLabel2("168L"), outerLabel2("169R"), outerLabel2("169L"),
			outerLabel("170R"), outerLabel("170L"), outerLabel("171R"), outerLabel("171L"),
			outerLabel2("172R"), outerLabel2("172L"), outerLabel2("173R"), outerLabel2("173L"),
			outerLabel("174R"), outerLabel("174L"), outerLabel("175R"), outerLabel("175L"),
			outerLabel2("176R"), outerLabel2("176L"), outerLabel2("177R"), outerLabel2("177L"),
			outerLabel("178R"), outerLabel("178L"), outerLabel("179R"), outerLabel("179L"),
			outerLabel2("180R"), outerLabel2("180L"), outerLabel2("181R"), outerLabel2("181L"),
			outerLabel("182R"), outerLabel("182L"), outerLabel("183R"), outerLabel("183L"),
			outerLabel2("184R"), outerLabel2("184L"), outerLabel2("185R"), outerLabel2("185L"),
			outerLabel("186R"), outerLabel("186L"), outerLabel("187R"), outerLabel("187L"),
			outerLabel2("188R"), outerLabel2("188L"), outerLabel2("189R"), outerLabel2("189L"),
			outerLabel("190R"), outerLabel("190L"), outerLabel("191R"), outerLabel("191L"),
			outerLabel2("192R"), outerLabel2("192L"), outerLabel2("193R"), outerLabel2("193L"),
			outerLabel("194R"), outerLabel("194L"), outerLabel("195R"), outerLabel("195L"),
			outerLabel2("196R"), outerLabel2("196L"), outerLabel2("197R"), outerLabel2("197L"),
			outerLabel("198R"), outerLabel("198L"), outerLabel("199R"), outerLabel("199L"),
			outerLabel2("200R"), outerLabel2("200L"), outerLabel2("201R"), outerLabel2("201L"),
			outerLabel("202R"), outerLabel("202L"), outerLabel("203R"), outerLabel("203L"),
			outerLabel2("204R"), outerLabel2("204L"), outerLabel2("205R"), outerLabel2("205L"),
			outerLabel("206R"), outerLabel("206L"), outerLabel("207R"), outerLabel("207L"),
			outerLabel2("208R"), outerLabel2("208L"), outerLabel2("209R"), outerLabel2("209L"),
			outerLabel("210R"), outerLabel("210L"), outerLabel("211R"), outerLabel("211L"),
			outerLabel2("212R"), outerLabel2("212L"), outerLabel2("213R"), outerLabel2("213L"),
			outerLabel("214R"), outerLabel("214L"), outerLabel("215R"), outerLabel("215L"),
			outerLabel2("216R"), outerLabel2("216L"), outerLabel2("217R"), outerLabel2("217L"),
			outerLabel("218R"), outerLabel("218L"), outerLabel("219R"), outerLabel("219L"),
			outerLabel2("220R"), outerLabel2("220L"), outerLabel2("221R"), outerLabel2("221L"),
			outerLabel("222R"), outerLabel("222L"), outerLabel("223R"), outerLabel("223L"),
			outerLabel2("224R"), outerLabel2("224L"), outerLabel2("225R"), outerLabel2("225L"),
			outerLabel("226R"), outerLabel("226L"), outerLabel("227R"), outerLabel("227L"),
			outerLabel2("228R"), outerLabel2("228L"), outerLabel2("229R"), outerLabel2("229L"),
			outerLabel("230R"), outerLabel("230L"), outerLabel("231R"), outerLabel("231L"),
			outerLabel2("232R"), outerLabel2("232L"), outerLabel2("233R"), outerLabel2("233L"),
			outerLabel("234R"), outerLabel("234L"), outerLabel("235R"), outerLabel("235L"),
			outerLabel2("236R"), outerLabel2("236L"), outerLabel2("237R"), outerLabel2("237L"),
			outerLabel("238R"), outerLabel("238L"), outerLabel("239R"), outerLabel("239L"),
			outerLabel2("240R"), outerLabel2("240L"), outerLabel2("241R"), outerLabel2("241L"),
			outerLabel("242R"), outerLabel("242L"), outerLabel("243R"), outerLabel("243L"),
			outerLabel2("244R"), outerLabel2("244L"), outerLabel2("245R"), outerLabel2("245L"),
			outerLabel("246R"), outerLabel("246L"), outerLabel("247R"), outerLabel("247L"),
			outerLabel2("248R"), outerLabel2("248L"), outerLabel2("249R"), outerLabel2("249L"),
			outerLabel("250R"), outerLabel("250L"), outerLabel("251R"), outerLabel("251L"),
			outerLabel2("252R"), outerLabel2("252L"), outerLabel2("253R"), outerLabel2("253L"),
			outerLabel("254R"), outerLabel("254L"), outerLabel("255R"), outerLabel("255L"),
			outerLabel2("256R"), outerLabel2("256L"), outerLabel2("257R"), outerLabel2("257L"),
			outerLabel("258R"), outerLabel("258L"), outerLabel("259R"), outerLabel("259L"),
			outerLabel2("260R"), outerLabel2("260L"), outerLabel2("261R"), outerLabel2("261L"),
			outerLabel("262R"), outerLabel("262L"), outerLabel("263R"), outerLabel("263L"),
			outerLabel2("264R"), outerLabel2("264L"), outerLabel2("265R"), outerLabel2("265L"),
			outerLabel("266R"), outerLabel("266L"), outerLabel("267R"), outerLabel("267L"),
			outerLabel2("268R"), outerLabel2("268L"), outerLabel2("269R"), outerLabel2("269L"),
			outerLabel("270R"), outerLabel("270L"), outerLabel("271R"), outerLabel("271L"),
			outerLabel2("272R"), outerLabel2("272L"), outerLabel2("273R"), outerLabel2("273L"),
			outerLabel("274R"), outerLabel("274L"), outerLabel("275R"), outerLabel("275L"),
			outerLabel2("276R"), outerLabel2("276L"), outerLabel2("277R"), outerLabel2("277L"),
			outerLabel("278R"), outerLabel("278L"), outerLabel("279R"), outerLabel("279L"),
			outerLabel2("280R"), outerLabel2("280L"), outerLabel2("281R"), outerLabel2("281L"),
			outerLabel("282R"), outerLabel("282L"), outerLabel("283R"), outerLabel("283L"),
			outerLabel2("284R"), outerLabel2("284L"), outerLabel2("285R"), outerLabel2("285L"),
			outerLabel("286R"), outerLabel("286L"), outerLabel("287R"), outerLabel("287L"),
			outerLabel2("288R"), outerLabel2("288L"), outerLabel2("289R"), outerLabel2("289L"),
			outerLabel("290R"), outerLabel("290L"), outerLabel("291R"), outerLabel("291L"),
			outerLabel2("292R"), outerLabel2("292L"), outerLabel2("293R"), outerLabel2("293L"),
			outerLabel("294R"), outerLabel("294L"), outerLabel("295R"), outerLabel("295L"),
			outerLabel2("296R"), outerLabel2("296L"), outerLabel2("297R"), outerLabel2("297L"),
			outerLabel("298R"), outerLabel("298L"), outerLabel("299R"), outerLabel("299L"),
			outerLabel2("300R"), outerLabel2("300L"), outerLabel2("301R"), outerLabel2("301L"),
			outerLabel("302R"), outerLabel("302L"), outerLabel("303R"), outerLabel("303L"),
			outerLabel2("304R"), outerLabel2("304L"), outerLabel2("305R"), outerLabel2("305L"),
			outerLabel("306R"), outerLabel("306L"), outerLabel("307R"), outerLabel("307L"),
			outerLabel2("308R"), outerLabel2("308L"), outerLabel2("309R"), outerLabel2("309L"),
			outerLabel("310R"), outerLabel("310L"), outerLabel("311R"), outerLabel("311L"),
			outerLabel2("312R"), outerLabel2("312L"), outerLabel2("313R"), outerLabel2("313L"),
			outerLabel("314R"), outerLabel("314L"), outerLabel("315R"), outerLabel("315L"),
			outerLabel2("316R"), outerLabel2("316L"), outerLabel2("317R"), outerLabel2("317L"),
			outerLabel("318R"), outerLabel("318L"), outerLabel("319R"), outerLabel("319L"),
			outerLabel2("320R"), outerLabel2("320L"), outerLabel2("321R"), outerLabel2("321L"),
			outerLabel("322R"), outerLabel("322L"), outerLabel("323R"), outerLabel("323L"),
			outerLabel2("324R"), outerLabel2("324L"), outerLabel2("325R"), outerLabel2("325L"),
			outerLabel("326R"), outerLabel("326L"), outerLabel("327R"), outerLabel("327L"),
			outerLabel2("328R"), outerLabel2("328L"), outerLabel2("329R"), outerLabel2("329L"),
			outerLabel("330R"), outerLabel("330L"), outerLabel2("331R"), outerLabel2("331L"),
			outerLabel2("332R"), outerLabel2("332L"), outerLabel2("333R"), outerLabel2("333L"),
		)
	}

	if err := g.WriteGerber(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done.")

	if *view {
		viewer.Gerber(g, false)
	}
}

func genPt(xScale, angle, halfTW, offset float64) Pt {
	r := (angle + *trace + *gap) / (3 * math.Pi)
	x := (r + halfTW) * math.Cos(angle+offset)
	y := (r + halfTW) * math.Sin(angle+offset)
	return Point(x*xScale, y)
}

type spiral struct {
	startAngle float64
	endAngle   float64
	size       float64
}

func newSpiral() *spiral {
	startAngle := 106.0 * math.Pi
	endAngle := 106.0*math.Pi + float64(*n)*2.0*math.Pi
	p1 := genPt(1.0, endAngle, *trace*0.5, 0)
	size := 2 * p1.Length()
	p2 := genPt(1.0, endAngle, *trace*0.5, math.Pi)
	if v := 2 * p2.Length(); v > size {
		size = v
	}
	return &spiral{
		startAngle: startAngle,
		endAngle:   endAngle,
		size:       size,
	}
}

func (s *spiral) genSpiral(xScale, offset, trimY float64) (startPt Pt, pts []Pt, endPt Pt) {
	halfTW := *trace * 0.5
	var endAngle float64
	if xScale < 0 { // odd
		endAngle = s.endAngle + 3.0*math.Pi/nlayers
	} else { // even
		endAngle = s.endAngle - 1.0*math.Pi/nlayers
	}
	// if trimY < 0 { // Only for layer2SpiralL - extend another Pi/2
	// 	endAngle += 0.5 * math.Pi
	// }
	steps := int(0.5 + (endAngle-s.startAngle) / *step)
	for i := 0; i < steps; i++ {
		angle := s.startAngle + *step*float64(i)
		if i == 0 {
			startPt = genPt(xScale, angle, 0, offset)
		}
		pts = append(pts, genPt(xScale, angle, halfTW, offset))
	}
	var trimYsteps int
	if trimY > 0 {
		trimYsteps++
		for {
			if pts[len(pts)-trimYsteps][1] > trimY {
				break
			}
			trimYsteps++
		}
		lastStep := len(pts) - trimYsteps
		trimYsteps--
		pts = pts[0 : lastStep+1]
		pts = append(pts, Pt{pts[lastStep][0], trimY})
		angle := s.startAngle + *step*float64(steps-1-trimYsteps)
		eX := genPt(xScale, angle, 0, offset)
		endPt = Pt{eX[0], trimY}
		nX := genPt(xScale, angle, -halfTW, offset)
		pts = append(pts, Pt{nX[0], trimY})
	} else if trimY < 0 { // Only for layer2SpiralL
		trimYsteps++
		for {
			if pts[len(pts)-trimYsteps][1] < trimY {
				break
			}
			trimYsteps++
		}
		lastStep := len(pts) - trimYsteps
		trimYsteps--
		pts = pts[0 : lastStep+1]
		pts = append(pts, Pt{pts[lastStep][0], trimY})
		angle := s.startAngle + *step*float64(steps-1-trimYsteps)
		eX := genPt(xScale, angle, 0, offset)
		endPt = Pt{eX[0], trimY}
		nX := genPt(xScale, angle, -halfTW, offset)
		pts = append(pts, Pt{nX[0], trimY})
	} else {
		pts = append(pts, genPt(xScale, endAngle, halfTW, offset))
		endPt = genPt(xScale, endAngle, 0, offset)
		pts = append(pts, genPt(xScale, endAngle, -halfTW, offset))
	}
	for i := steps - 1 - trimYsteps; i >= 0; i-- {
		angle := s.startAngle + *step*float64(i)
		pts = append(pts, genPt(xScale, angle, -halfTW, offset))
	}
	pts = append(pts, pts[0])
	return startPt, pts, endPt
}
