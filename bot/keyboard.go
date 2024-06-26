package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var carMakes = []string{
	"Audi",
	"BMW",
	"Ford",
	"Jaguar",
	"Land Rover",
	"Mercedes",
	"Nissan",
	"Porsche",
	"Toyota",
	"Vauxhall",
	"Volkswagen",
	"Volvo",
}

var AllCarMakes = map[string]struct{}{

	"Abarth":                     {},
	"AC":                         {},
	"Aixam":                      {},
	"AK":                         {},
	"Alfa Romeo":                 {},
	"Alpine":                     {},
	"Alvis":                      {},
	"Ariel":                      {},
	"Aston Martin":               {},
	"Auburn":                     {},
	"Audi":                       {},
	"Austin":                     {},
	"BAC":                        {},
	"Beauford":                   {},
	"Bentley":                    {},
	"BMW":                        {},
	"Bugatti":                    {},
	"Buick":                      {},
	"BYD":                        {},
	"Cadillac":                   {},
	"Caterham":                   {},
	"Chesil":                     {},
	"Chevrolet":                  {},
	"Chrysler":                   {},
	"Citroen":                    {},
	"Corvette":                   {},
	"CUPRA":                      {},
	"Dacia":                      {},
	"Daewoo":                     {},
	"Daihatsu":                   {},
	"Daimler":                    {},
	"Dax":                        {},
	"De Tomaso":                  {},
	"DFSK":                       {},
	"Dodge":                      {},
	"DS AUTOMOBILES":             {},
	"Ferrari":                    {},
	"Fiat":                       {},
	"Fisker":                     {},
	"Ford":                       {},
	"Gardner Douglas":            {},
	"Genesis":                    {},
	"GMC":                        {},
	"Great Wall":                 {},
	"GWM ORA":                    {},
	"Hillman":                    {},
	"Honda":                      {},
	"Hummer":                     {},
	"Hyundai":                    {},
	"INEOS":                      {},
	"Infiniti":                   {},
	"Isuzu":                      {},
	"Iveco":                      {},
	"Jaguar":                     {},
	"JBA":                        {},
	"Jeep":                       {},
	"Jensen":                     {},
	"KGM":                        {},
	"Kia":                        {},
	"KTM":                        {},
	"Lada":                       {},
	"Lagonda":                    {},
	"Lamborghini":                {},
	"Lancia":                     {},
	"Land Rover":                 {},
	"LEVC":                       {},
	"Lexus":                      {},
	"Leyland":                    {},
	"Lincoln":                    {},
	"Lister":                     {},
	"London Taxis International": {},
	"Lotus":                      {},
	"Mahindra":                   {},
	"Marcos":                     {},
	"Maserati":                   {},
	"MAXUS":                      {},
	"Maybach":                    {},
	"Mazda":                      {},
	"McLaren":                    {},
	"Mercedes-Benz":              {},
	"Merlin":                     {},
	"MEV":                        {},
	"MG":                         {},
	"Microcar":                   {},
	"MINI":                       {},
	"Mitsubishi":                 {},
	"MK":                         {},
	"MOKE":                       {},
	"Morgan":                     {},
	"Morris":                     {},
	"Nardini":                    {},
	"NG":                         {},
	"Nissan":                     {},
	"Noble":                      {},
	"Opel":                       {},
	"Pagani":                     {},
	"Panther":                    {},
	"Perodua":                    {},
	"Peugeot":                    {},
	"PGO":                        {},
	"Polaris":                    {},
	"Polestar":                   {},
	"Pontiac":                    {},
	"Porsche":                    {},
	"Proton":                     {},
	"Radical":                    {},
	"Reliant":                    {},
	"Renault":                    {},
	"Replica":                    {},
	"Riley":                      {},
	"Robin Hood":                 {},
	"Rolls-Royce":                {},
	"Rover":                      {},
	"Saab":                       {},
	"SEAT":                       {},
	"Skoda":                      {},
	"Smart":                      {},
	"SsangYong":                  {},
	"Stuart Taylor":              {},
	"Subaru":                     {},
	"Sunbeam":                    {},
	"Suzuki":                     {},
	"Tesla":                      {},
	"Tiger":                      {},
	"Toyota":                     {},
	"Triumph":                    {},
	"TVR":                        {},
	"Ultima":                     {},
	"Vauxhall":                   {},
	"Volkswagen":                 {},
	"Volvo":                      {},
	"Westfield":                  {},
	"Yamaha":                     {},
}

var KeyboardMap = map[string]tgbotapi.InlineKeyboardMarkup{
	"popular": tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Audi", "Audi"),
			tgbotapi.NewInlineKeyboardButtonData("BMW", "BMW"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Ford", "Ford"),
			tgbotapi.NewInlineKeyboardButtonData("Jaguar", "Jaguar"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Land Rover", "Land_Rover"),
			tgbotapi.NewInlineKeyboardButtonData("Mercedes", "Mercedes"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Back", "begin"),
			tgbotapi.NewInlineKeyboardButtonData("Cancel", "cancel"),
			tgbotapi.NewInlineKeyboardButtonData("Next", "nextPopular"),
		),
	),
	"nextPopular": tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Nissan", "Nissan"),
			tgbotapi.NewInlineKeyboardButtonData("Porsche", "Porsche"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Toyota", "Toyota"),
			tgbotapi.NewInlineKeyboardButtonData("Vauxhall", "Vauxhall"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Volkswagen", "Volkswagen"),
			tgbotapi.NewInlineKeyboardButtonData("Volvo", "Volvo"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Back", "popular"),
			tgbotapi.NewInlineKeyboardButtonData("Cancel", "cancel"),
		),
	),
	"begin": tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Popular Makes", "popular"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("!German", "german"),
			tgbotapi.NewInlineKeyboardButtonData("!British", "british"),
			tgbotapi.NewInlineKeyboardButtonData("!Japanese", "japanese"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("!Other", "other"),
			tgbotapi.NewInlineKeyboardButtonData("Cancel", "cancel"),
		),
	),
	"model": tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Models", "back"),
		),
	),
}

func createKeyboard(r []string) tgbotapi.InlineKeyboardMarkup {

	var inlinekeys [][]tgbotapi.InlineKeyboardButton

	for x := 0; x < len(r)-1; x += 2 {
		var row []tgbotapi.InlineKeyboardButton
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(r[x], r[x]))
		inlinekeys = append(inlinekeys, row)
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(r[x+1], r[x+1]))
		inlinekeys = append(inlinekeys, row)

	}
	return tgbotapi.NewInlineKeyboardMarkup(inlinekeys...)
}
