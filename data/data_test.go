package data_test

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"khan/kutil/container/kslices"
	data_utils "khan/kutil/data"
	"khan/kutil/filesystem"
)

// publish aa "\x78\x01\x73\x09\xF1\x88\xF4\xF2\xB6\x72\x72\x77\xF6\x70\x77\xB3\x0A\x34\xB4\xF2\x34\x30\x30\xB4\x8A\x70\x02\x91\x7E\x4E\x81\x20\xCA\xC3\x27\x02\x44\x39\x67\x80\x48\x17\x67\x47\x07\x37\x2B\x43\x13\x3D\x73\x2B\x43\x33\x0B\x53\x13\x43\x13\x73\x4B\x4B\x00\x5F\x2B\x10\x94"
func TestUnZip(t *testing.T) {

	str := "78017309F188F4F2B6727277F67077B30A34B4F2343030B48A7002917E4E8120CAC32702443967804817674707372B43133D732B43330B53134313734B4B005F2B1094"
	// str = "789c9d96793c13fe1fc7e74eb9ca2d9411df2421e6ce99db900ce5c8554deefb0a432464cefc98f3ebbe6f8d45e418332624b3213315ad62b16cc94f8fdf1fbffffbe7f3f87c1e9fff9ecfcfebfdfa0000004e0000a008f0077800a08000803ec018701b6006300458002c015a00a7939329c0e4646707503e596f9fdc5f0700986e44216a1686034694b8d8d785f10a2bc11b32754c8277d5b9990c9f946a7ce95c51440d15e6332d56a7576391ee12da936b5fdcb6da0277db63ddc6b9c78444b39b0ed9d349fd8f1c891dcecdb4953a1f99e4f552a0858c704a8a1db9f2b0fcdba00da32d955fcd9b49b3dac128d41ae2c9bc2f18042358d6255d9c56bbb4763ddc4ff7e1a94193f8a0a3086cf888c4ad0d371d56e3cd4e4edee856c9a4a44cd27db11784d8aa45dc56ffd46b4efb65736e64c978cf4776bffea2ed4664f9f89b686acf73256d1be912d9e2f1c3e9a652d1cfd5911a0168ff0493cb71ef85a6963f7ece8d60cab469eb5060b322b86d4a333fe33c921f4bc885dc4a7dbfd021b467f0b6932d2fc2ba86db3d0f94f0d9c0051a46cc85a0977fb0635f117e52427a8b7afe03a54f4a04ed60ceca2b49af0d182ca53a4968aaa394f5a6a67baefc15f3eb7f985bdb301b9c30d7abb65ea1c485d99c3e7f4feddc45d61fa306bc312aca885ffb33338952a3a397a78fccd7a00ff02d0ffe4d9b3888d62dedb6f48df55ec59a5acaf80997c40f2e879324c9016f73015db5a43b7bf232990601fd625f7951a251aad4a075dfdcf9b6ecbee5c8aaaa0aa6b3e9578e8e76e23b5c2f4d8a0cc7daf9cb93b7ff093c8bf3d01c8efa60e4ef1964fcb4a761a99c41899d48d7920d1ab6d59caf9c1477e15a6f8f30e17c9e93e12d70cad382910f9d325f179d34747f5ff622fe2ecac0c2c52b450965a67da1b14925e868408b8c0d27f804d850ae85556a28fdaab2976c23152339ac116368860e147cec944a266111983dbbe6f7aed61fbe3687a7fc98dda6b38c5b560aa16b88107011b29aec551ba64996c61a81a85b5a8a4ff8d1a2c74c4786c3fc7f05dae77fa0dffc79dc7aac85140df352febb988a8be3206b8606255135380933fa4d0ebae87f51f552f3c71b9d47b2c15d71e27aba69d83d8f735b0277bcb930b34647539179af57c5eca5a2cea31be04a3148e86fd91d5ab88adb4d41b4d68ccab5493ee9582c3d618546466b9e89435e7b4b0cd3316d1ef15b043b0ed496b508a7e7315d82140057e3e388df68738fd8fcf8e636e9ad5455b98e700152f39e4716c1f23ed30f67db6c96be87b9bdfac583a87e3193c03429797647f34117ae3d21264655636bfb4e652f0bc91387348a5456bfac78e5546a07bcb69aabd2b61b492d03a32bf2a29471f9f95fdaf43ab5e045104e6248f050cac453c787035337368aecdf19d18789e1049af8f164d3ce4b1eec0ac13478a6414fb58ed824df0a3e01cdf2428341beca28c9deb62be916f99e521cf961e647140e721528b8ddbfc42c7323e1fbbca221af60e86f9844a49dcbcdd5881274795b5560611eb2fdafd4989fa859fcffdcd1ce22beace3d6ea0d1063be3d1cb1a1c449e976c59525045cae3f67f65295663394f7dedc79c8ca6ded59363d1aa2d050afe4bd72e89d3e3ae6e98a5d5a09b923a53bfdec22bc12bf91642c3d96188dbaa2aba30e3eeb848384ed4db41d7976f926d9e96b00b7c048ada958d68d8585e6d505af395755a94dcf39285f5bfcefb2086d6a3160f9009133bf34c4c0468a66aaa86c62a379abacef4b5e4abcc9ae87f07c58c3d44e9982a369859d4ec8c6d5505ff2834f5397ccac066d3589b87e8a0590ebf5b6f5d7d0000d95643cf7214fd3ac8ff8ecc4aeeaf31d98dc79a8f19a3e4de3038c74cfac6f046760a62e88c8227590950d50f1f88590cd2e7017e7f04e560c0ca311ff4935bdec95fb3a547ca9354a4ee715bd9a59993f61f3f3629fa91ec83e78bf49b3e9a0a8f940926df1770295f772c89b82bfc26ff32719ba27c9e03b49460a38f263db80d5a6d193c47a7ca32f6601f624abe6d8fd3d67432e53fb4af3dcbfdaaf7ad7028e092814a4ee68be117f6b51fa7334d12fa71d2f4edfff40e829217d0d76503eeb24fdb417d856a8e0d5ce4b0c44b4a24c9f2314a2ea76f386a245a2c43b192f3393cfc2edd2ae79e3c95bd5a745b8b7f68e96c83f976b5a24d8d6bbc3bebe63b1868a2a17ea974ea651f8044ad69793555e64c659de7e2288a8ff156a6c499f2133fa1f8da6ed86790dae3f6b4096e22fdce12202929f445ded8e6b1dca4ecb39af54149f1f82ccde19753d3042048bd9e039701aaf34a255faa8fd2d8c6a6e664b7d85819af0152bbc4a01bb631fa57af1ae4379c34ec36b336678de2f2475fb8c17d51c4ab0df9712e167c1b8c6401effa30f8d0978d7b56d94f54dbaeb3ab18a6e650b879f5f68a6a3aa0fa60621195b193b33efd52d2487dafe4a42f089847c2cf844428afb29c8c0017e079b0eb25b399120effb6181455c7fc3a65eae4c92a318646bebff5df9ed0bd7d9cc3a5fd937d48c8d46d1ec679a826a62776cd2428d1da70bf7c390ade5205e2c09b35ccf9794b8eb248fbec27f5ee04643fcb10115e54c7490ad9a11e1335495e5c8aaf0f0fa27b1be3af952316fae1337df03ffd8e5aac65583f6a7521caf40870c7a0b3a39c4d6b45881256fd618c353038f58e291185a5d03af02a74fcbcf395b744d08279760fd7a1e4e98a4bbe2ea9d80c84f66ac405f524ab4a5a38ffa324ca53104b86a29ed7ccfceba6be9db88d9f0c477a3f5cc9c728fea752785a736050dc2261d422e021a0eb58f45aed30a065b784e1935beececdb17901a3f753db380c58ba55dbdf8674a75a0827c4aefbf2e3cfd05cb771fee6b25edd15415e77a2cbd82e04a871c8463de173cb45dc12dd2876f3f5bbbc59d9efb4757ebf07cf230cd08a6d07462566224ad9efec841ecf18ef9c08c87527f48b0270d65fc9526933f59c1fcd1c486bb903497a91bdf5a914507d8f538948c341afa305d91591074cf26ca58dc0491163066a683cfef07f4dcb52a2b23ed04dcc838ffcc5226d2c51af414d4ba436ebc01992ffeee21c499c5c13fc92e0bf9f4b355cc76cd5682564733d6915602ca0b756e8d1fce2e4d6ca847f3f0b9ec55a7fcde0d73bbd0abe9875b5f96d6aa24449c83e9a4c64d193ace7fbdcff14e5e73ed71396f79c1d46cb8d03f9bc70e86ee2adc78adefa7d50a707d569352a21c2622458c0efd50d133605ef09025486c89b03aeefae6f175cf31c2dab891b09b9c9bf820b13e4c3b2b5a7bcedfa50725b2dd011f78ede3310161a23954f00d57bd5094a139046527faed8678c78e1d4fe94b25f74a5d4364d27cee99b1144cbf7310fd09cffb98cdd67d35b922659a2be375a23007e80d5004b3026c87c7de0515df776235facd61193cf277ec03a310a5ff6ff01cca164677ae070fe19b0739482022125910701cf8570ca2b505d801e57c05b179a43db312b63346b23f7cdb5248e6c108995d50cf9048bc4a0baca3781cb72b7c338019f16493cf5d7db4c00a5e535bdaed0d543fea689875d7ef9d121e96dbfbba13f6f3a1593d539dd22fc2a62d6ce691d1e99aef614e2cdcb9b6db372d18358c7905b3b19b5139e3a4fb6744d5ba66ab66612a1df6eccb029a3fd46cc3c07721f72ea7d3916bf122c249bb53b2e5b6d3195c60165344f419ada096d3fc10ea385a3128764b23a81c7b845643877a6d0b6b64eb3cc439a31a9b227026504b648be3bee7db4ec5b731b5ad7214ce90fd6d06f5f191b45b97a75c1629abb4e4b15acd5998e3f3707a3a5653762997a037aeebea67987b1b7eed1e305b1194a3c20ac753043c6fb53fb8ad183db4c5c70ff5bea52f1a9279ced97c327753eaa641aee3d59be689306934447f3bdf3320b5d3474eea507e82e0461a8a3b1ac9fd2b494e7fbadcf28f24c0499900578225bfa62b2efef96639ec44e061181e140fe3180835682a52cdcd79441d037bd65947eac54635f71289a83b1e8bccf423cf5e66a5c59079451043d64f3fe0a320b3768dd463fd2b9956ebd3d6a13343ae1286252293e5150c8d83b4d40dbf5e1ec1f3150ab9854b908cb9d3a1c651b530915d21aa59c2a7bc0e7d1e81aecdffd479b31a377bb647df37ada3b8b0b2be61ab2c1bed9cc903ca812d3e89a860f79ff94b92ed9d2f72f978b8288f3e92cac63562d4f64bde554afadac5b27d21f5786fd3c92fbf2c993fb3cfe4b9a84edecf667de22993d5995fa214c36cabc8ceaeb84161ecf6c43232cfbccc37105fa6697f3649eff0123c5388cbefae7defd57040f50088c958103a720814301d2b254f1a139e9328f8a8fd586ce296d0edcf8838657596beb04b7dc7a4804dd9aed2da05779b8921f31858d7c0086248eb77660253d6db36bfff02570f578d"
	data, err := hex.DecodeString(str)
	if err != nil {
		fmt.Println(err)
		t.Errorf("%s", err.Error())
		return
	}

	u8Data := []uint8(data)
	buf, err := data_utils.UnZip(u8Data)
	if err != nil {
		fmt.Println(err)
		t.Errorf("%s", err.Error())
		return
	}
	orgStr := string(buf)

	fmt.Printf("origin: %s\n", orgStr)
}

func TestZip(t *testing.T) {
	str := "DTHYJK:BGCHGF:Q1:I001:XB001:NBQ001:HLX001:Ch001:DCA@F:14.7:1685414799"
	u8Data := []uint8(str)
	buf, err := data_utils.Zip(u8Data)
	if err != nil {
		fmt.Println(err)
		t.Errorf("%s", err.Error())
		return
	}
	// orgStr := string(buf)

	fmt.Printf("compress data: %s\n", hex.EncodeToString(buf))

	// buff := []byte{120, 156, 202, 72, 205, 201, 201, 215, 81, 40, 207, 47, 202, 73, 225, 2, 4, 0, 0, 255, 255, 33, 231, 4, 147}
	// b := bytes.NewReader(buff)
	// r, err := zlib.NewReader(b)
	// if err != nil {
	// 	panic(err)
	// }
	// io.Copy(os.Stdout, r)
	// r.Close()

	// var in bytes.Buffer
	// w := zlib.NewWriter(&in)
	// w.Write([]byte("hello, world\n"))
	// w.Close()

	// // zip := DoZlibCompress([]byte("hello, world\n"))
	// // fmt.Println(zip)
	// fmt.Println(hex.EncodeToString(in.Bytes()))
}

func TestTrim(t *testing.T) {
	header := "'sync_data ;1.0 ;1686815770  ;1686815800  ;recv_center001          ;recv_center001          ;                                ' \r\n\t"
	header = strings.Trim(header, "' \r\n\t")
	header = strings.TrimSpace(header)
	slice := strings.Split(header, ";")

	kslices.Process[string](slice, func(val string) string {
		result := strings.TrimSpace(val)
		fmt.Printf("[%s]\n", result)
		return result
	})
	// fmt.Printf("[%s]\n", slice)
}

func TestFileReExtName(t *testing.T) {
	filePath := "C:\\Private\\Project\\Golang\\SFtpK\\logs\\SFtpKc.2023-07-18.log"
	result := filesystem.FileReExtName(filePath, ".txt")
	if result == false {
		t.Errorf("%s", "file not found")
	} else {
		fmt.Printf("ok\n")
	}

	// filePath = "C:\\Private\\Project\\Golang\\SFtpK\\logs\\SFtpKc.2023-07-18"
	// result = filesystem.FileReExtName(filePath, ".txt")
	// if result == "" {
	// 	t.Errorf("%s", "file not found")
	// } else {
	// 	fmt.Printf("%s \n", filesystem.FileReExtName(filePath, ".txt"))
	// }

	// filePath = "C:\\Private\\Project\\Golang\\SFtpK\\logs\\SFtpKc.2023-07-18."
	// result = filesystem.FileReExtName(filePath, ".txt")
	// if result == "" {
	// 	t.Errorf("%s", "file not found")
	// } else {
	// 	fmt.Printf("%s \n", filesystem.FileReExtName(filePath, ".txt"))
	// }

	// filePath = "C:\\Private\\Project\\Golang\\SFtpK\\logs\\SFtpKc_2023-07-18..."
	// result = filesystem.FileReExtName(filePath, ".txt")
	// if result == "" {
	// 	t.Errorf("%s", "file not found")
	// } else {
	// 	fmt.Printf("%s \n", filesystem.FileReExtName(filePath, ".txt"))
	// }
}

func TestTimestamp(t *testing.T) {
	str1 := "1693474484000"

	timestamp, err := strconv.ParseInt(str1, 10, 64)
	if nil != err {
		timestamp = 0
	}

	fmt.Printf("timestamp sec: %s\n", time.Unix(0, timestamp*1000*1000).Format("2006-01-02 15:04:05.000"))

	str2 := "1693474484"

	timestamp, err = strconv.ParseInt(str2, 10, 64)
	if nil != err {
		timestamp = 0
	}

	fmt.Printf("timestamp sec: %s\n", time.Unix(timestamp, 0).Format("2006-01-02 15:04:05.000"))
}

func Test_TrimStr(t *testing.T) {
	b := []byte{116, 101, 115, 116, 48, 48, 49, 0, 0, 0, 0, 0}
	tmp := string(b)
	// str := strings.TrimSpace(tmp) //不能去除 \x00
	str := strings.Trim(tmp, "\x00 \b\t\n\r")

	fmt.Printf("%s\n", str)
	fmt.Printf("%v\n", []byte(str))
}
