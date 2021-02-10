package mobile

import (
	"crypto/rsa"
	"math/big"
)

var (
	rsaE = 0x10001

	rsaClientN = "8C81424BC166F4918756E9F7B22EFAA03479B081E61896872CB7C51C910D7EC1" +
		"A4CE2871424D5C9149BD5E08A25959A19AD3C981E6512EFDAB2BB8DA3F1E315C" +
		"294BD117A9FB9D8CE8E633B4962E087C629DC6CA3A149214B4091EF2B0363CB3" +
		"AE6C7EE702377F055ED3CD93F6C342256A76554BBEA7F203437BBE65F2DA2741"

	rsaClientD = "3704DAB00D80C25E464FFB785A16D95F688D0A5823811758C16308D5A1DB55FA" +
		"800D967A9B4AEDE79AA783ADFFDCDB23541C80B8D436901F172B1CCCA190B224" +
		"DBE777BF18B96DD9A30AACE8780350793A4F90A645A7747EF695622EADBE23A4" +
		"C6D88F22E87842B43B35486C2D1B5B1FA77DB3528B0910CA84EDB7A46AFDBED1"

	rsaClientP = "E98CA62CB05FA4ABD8B203D5AF5D18DB63DF9C6B6ED87DAAC3F56573592DAB16" +
		"0F0CB026EF0A8F5E6D77268EEE384210A0850148557B9E6D0ED0A7276FA85D25"

	rsaClientQ = "9A02E789BD57A2760EC635493368EACB9CC419EEAFCE0F1A4B028261E735E228" +
		"892A611870FE330D2466B38DE19D0B29F0CDB29E39EC6028F289E820F8067CED"

	rsaServerN = "D1DEB5A67989D10A5C0B041D11AB47659BC32B3E65A92F9236DD2CC2140A0DA4" +
		"95FA0E7EF875F6FACB4D981F5E5E3BA5E67811DD73B7A4105254222FACC206A7" +
		"4E8911215451AFFCF5EBFDE887BC9ADB675C02A835F75CBE77A603AF820683F4" +
		"39E7CED4289323DF9C9055769E15798B9453E89058B3E2E3F186B07B046A996D"

	//rsaClientDp = "36EC309809D23433857E3790A4F0CBCBAC2D05E7EDE553883915188A8BCA4595" +
	//	"A66C6170867E8140BF9569A7EB35A7B3A94C1E0518B53D8880176977C8B65B51"
	//
	//rsaClientDq = "7FF9CDE07CFFA73626CCB9569C6BA03F9582B671CA9095A829906A3B645F3810" +
	//	"AAFA1638B31BE7DC11D56D7A867172E764FBE862E68AEED4D7C594A860B13379"
	//
	//rsaClientQp = "BF0CB7702417519C8B91CEBA0EA34F8B3867EA830EB517A3D0654F04963FDD2F" +
	//	"C251C5FA691326929DE5669C79FD2F57EC6E9AA368AEF2E921758AC3135AA073"
)

func (c *Client) rasInit() {
	// TODO: Change to pem data?

	// RSA private key
	c.rsaPrivKey = &rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: big.NewInt(0),
			E: rsaE,
		},
		D: big.NewInt(0),
		Primes: []*big.Int{
			big.NewInt(0), big.NewInt(0),
		},
	}
	c.rsaPrivKey.N.SetString(rsaClientN, 16)
	c.rsaPrivKey.D.SetString(rsaClientD, 16)
	c.rsaPrivKey.Primes[0].SetString(rsaClientP, 16)
	c.rsaPrivKey.Primes[1].SetString(rsaClientQ, 16)
	c.rsaPrivKey.Precompute()
	// RSA public key
	c.rsaPubKey = &rsa.PublicKey{
		N: big.NewInt(0),
		E: rsaE,
	}
	c.rsaPubKey.N.SetString(rsaServerN, 16)

}
