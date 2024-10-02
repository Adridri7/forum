package users

import (
	"math/rand"
)

func RandomProfilPicture() string {
	randomPP := []string{
		"https://www.mesopinions.com/public/img/petition/petition-img-153019-fr.jpeg",
		"https://koreus.cdn.li/media/201404/90-insolite-34.jpg",
		"https://www.wipo.int/export/sites/www/wipo_magazine/images/en/2018/2018_01_art_7_1_400.jpg",
		"https://i.etsystatic.com/40574730/r/il/cc7a16/4592422959/il_fullxfull.4592422959_7mbc.jpg",
		"https://ih1.redbubble.net/image.3546319896.8009/flat,750x,075,f-pad,750x1000,f8f8f8.jpg",
	}

	random := rand.Intn(len(randomPP))

	return randomPP[random]
}
