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
		"https://images3.memedroid.com/images/UPLOADED861/61db381c190fd.jpeg",
		"https://cdn.discordapp.com/attachments/1219314334778265811/1291071596160876604/hi-fwends.jpg?ex=66fec39f&is=66fd721f&hm=4d3ad49db64ba36d75efa6a26a699a3accb57dfb6a44d867f4ce78b5fc830fc3&",
		"https://cdn.discordapp.com/attachments/1219314334778265811/1291071754588131329/Screenshot_20231209-222748-378.png?ex=66fec3c5&is=66fd7245&hm=dcd05183a25b22914599e54cdbacf7cb21e29651f6fbd42aa03e129f268da00e&",
		"https://koreus.cdn.li/media/201812/disaster-macron.jpg",
		"https://cdn.discordapp.com/attachments/1219314334778265811/1291071936863932527/Capture_du_2020-04-28_01-59-17.png?ex=66fec3f0&is=66fd7270&hm=a385319b24c275e2247702f698048013be4e2f24767a4230c91a701a50cb55be&",
		"https://cdn.discordapp.com/attachments/1219314334778265811/1291072008943308883/image0-199.jpg?ex=66fec401&is=66fd7281&hm=e5c3bac7123d506b9bcc8b44cc0901e5fe68b687998932a21a42c65bc2073ef3&",
		"https://cdn.discordapp.com/attachments/1219314334778265811/1291072316142260396/FB_IMG_1688529985879.jpg?ex=66fec44b&is=66fd72cb&hm=ce37629e0594e4c893a306d4867865cfc9dab96b3b2e23030c3dd106ff0ebde4&",
		"https://i.redd.it/dye318nrz7c91.jpg",
		"https://cdn.discordapp.com/attachments/1219314334778265811/1291074702621544499/unknown-5_1661130584137.png?ex=66fec684&is=66fd7504&hm=78a83cd18fb3dddbedc6ba55e05d9088695783320804642fdd50dbbe13331356&",
		"https://cdn.discordapp.com/attachments/1219314334778265811/1291074601874358282/1657986172470.jpg?ex=66fec66c&is=66fd74ec&hm=c499849316c1a472f72a3e434482bcaadc7fb26501aa940a3e277dbec72b0fc4&",
		"https://cdn.discordapp.com/attachments/1219314334778265811/1291074780388266105/854084070157058089.png?ex=66fec696&is=66fd7516&hm=dc8c061bbd057997cdce7b5b90923880236c335d8da69ca5aaaf80ccd77577af&",
		"https://cdn.discordapp.com/attachments/1219314334778265811/1291074957266260131/FB54D9BB-4AC3-4379-80A4-937D2675D2C1.jpg?ex=66fec6c0&is=66fd7540&hm=a00976df6e2fea8d3d60bada24121f511f3814e3ff1c7e80bc7a64f7811112bb&",
		"https://cdn.discordapp.com/attachments/1219314334778265811/1291075058759897139/879696169524150282.webp?ex=66fec6d9&is=66fd7559&hm=942b9b2a5e25adab11e74978c914394752804639c054785996943116085d0203&",
		"https://cdn.discordapp.com/attachments/1219314334778265811/1291075152259321996/for_real.png?ex=66fec6ef&is=66fd756f&hm=614117874dd6ef7abb50236c40a76d6011c6bb8385fdcd5234f215bbff0dc91e&",
		"https://cdn.discordapp.com/attachments/1219314334778265811/1291075284984008705/E7kJ2YyVkAATEPR.png?ex=66fec70e&is=66fd758e&hm=c4389b072720bf7db95eca7dad2b98fe871330abec22a0e7b60a5d97db5e4c29&",
		"https://cdn.discordapp.com/attachments/1219314334778265811/1291075414068039824/E89BQdlX0AUoOm0.png?ex=66fec72d&is=66fd75ad&hm=90bfde9e464cd6c6a4f24762da74eecd92d7dbdb21ad4a90cd32df0df273ffa2&",
		"https://cdn.discordapp.com/attachments/1219314334778265811/1291075858957602879/E2YVvO5VkAAW9v4.png?ex=66fec797&is=66fd7617&hm=4651e0c0221f81819083e4b3a9035aa07135754941041a1254530244fb0a0d4d&",
		"https://cdn.discordapp.com/attachments/1219314334778265811/1291076172159123569/minecraff_baby.png?ex=66fec7e2&is=66fd7662&hm=b6c2682288d5323303c0bd643daca1f530db6e3903c2fae5a7e1b0af8583126b&",
		"https://cdn.discordapp.com/attachments/1219314334778265811/1291077125541072908/A.png?ex=66fec8c5&is=66fd7745&hm=1b7ebc971e981e51ebeda74c32864670736078c2011108a1c6cb70cafe3a754f&",
		"https://cdn.discordapp.com/attachments/1219314334778265811/1291077426947948574/AUUHGGHEEU.jpg?ex=66fec90d&is=66fd778d&hm=a0b8d7fc3feb5d46018a2d74dd1fdb59bcf7f9ca6d234805f3f99236a72b6dff&",
		"https://cdn.discordapp.com/attachments/1219314334778265811/1291077526692692078/Capture_du_2019-12-21_12-34-19.png?ex=66fec925&is=66fd77a5&hm=96b87610d002344ff1b6faae79c03b8bfbb9d5d7aecf7326c4487f26dace7ad4&",
	}

	random := rand.Intn(len(randomPP))

	return randomPP[random]
}
