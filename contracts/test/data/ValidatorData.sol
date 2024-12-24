// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

/* solhint-disable max-line-length */

contract ValidatorData {
    struct ValidatorInfo {
        bytes uncompressedHex;
        bytes compressedHex;
        address evmAddress;
    }

    ValidatorInfo[] validators;

    constructor() {
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"0481513466154dfc0e702d8325d7a45284bc395f60b813230da299fab640c1eb085a6e35cb30d41af355e4701f455447a675c497f98d38e783ba3f91bf75595a8f", // pragma: allowlist secret
                compressedHex: hex"0381513466154dfc0e702d8325d7a45284bc395f60b813230da299fab640c1eb08", // pragma: allowlist secret
                evmAddress: 0xb5cb887155446f69b5e4D11C30755108AC87e9cD
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"04dd68a0a4923a1b9321d39f01425f7b631066514cb2e5e1b5ed91e5c327d30c534b55b55f63ca2ba7335f8086b955d93c25225986d8d8992675ce526fe720ecf1", // pragma: allowlist secret
                compressedHex: hex"03dd68a0a4923a1b9321d39f01425f7b631066514cb2e5e1b5ed91e5c327d30c53", // pragma: allowlist secret
                evmAddress: 0xf89D606F67a267E9dbCc813c8169988aB8aAeB5E
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"04b95e0ff2299de8215fda840dcf37c9398296b0000773a7cf96e6814b979806acce72cb5239fd06f3d86c5ed409b3b5844182ceabb696835a9f9f2a9392ed2add", // pragma: allowlist secret
                compressedHex: hex"03b95e0ff2299de8215fda840dcf37c9398296b0000773a7cf96e6814b979806ac", // pragma: allowlist secret
                evmAddress: 0x7D42B46Cc15bC492C928950acD34E3A474305B8c
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"04a5e10cd81da0bfe7644673fb349dea42d6b058b126a6dc81526561cc9d1680aae5b7adfaccba3c3c48d2b2ed637e6bc4eaca65da45b4ac57278c16daccf2dbbe", // pragma: allowlist secret
                compressedHex: hex"02a5e10cd81da0bfe7644673fb349dea42d6b058b126a6dc81526561cc9d1680aa", // pragma: allowlist secret
                evmAddress: 0xED0A0fb9d98a8C1d4D6E4A47241B79eEE982A42c
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"04e40b0d98016dd26fdb19f446bc04739b3728f9e35ebf75b13584d745c9f360bc9adf8910018442f12fda19a323e321c6f7f68e68d13e27fb39edf667255a3efc", // pragma: allowlist secret
                compressedHex: hex"02e40b0d98016dd26fdb19f446bc04739b3728f9e35ebf75b13584d745c9f360bc", // pragma: allowlist secret
                evmAddress: 0x47F4fdC22191BB622a9Ea2645E684e12715D2512
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"048400994469b03d92a8c81949125c799a6d3ef26e4c16e200aec20a3f9fdd6a00ac491c3192fe316d71d866ed102c59c4fc0073f06226202e34e88ccec9bb8f9d", // pragma: allowlist secret
                compressedHex: hex"038400994469b03d92a8c81949125c799a6d3ef26e4c16e200aec20a3f9fdd6a00", // pragma: allowlist secret
                evmAddress: 0x9f6f533a1620f5c66539dF2dF3a2Ff3E93a36830
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"049d511cb4290718656a8eab48475d6127fc33acfc3e6f0415f5c4642d7766db554faef7531f593eff2c1c0eeb29881c094055a91a158526484a95c2f08e598d83", // pragma: allowlist secret
                compressedHex: hex"039d511cb4290718656a8eab48475d6127fc33acfc3e6f0415f5c4642d7766db55", // pragma: allowlist secret
                evmAddress: 0x54DA9dC7A141E86f6D5b8bb3Bd770f8d6F606DA7
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"046d6057a9982310f0d907edc08079912e09ca4dc843c80b26c1107262f11e1580b81db4439e93324309271730135a885d90b01e3a134da12bed39868cb4e9f94a", // pragma: allowlist secret
                compressedHex: hex"026d6057a9982310f0d907edc08079912e09ca4dc843c80b26c1107262f11e1580", // pragma: allowlist secret
                evmAddress: 0x53b9e8B9C04a6f6d4ae7BdAC759dE986083e7F0f
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"04bfaa0a7de99fcca7702e00abff8bc9be4191bdd3b91d6900c4076c4bfdeb21f65793df6d75f88a7f5db5af53384c23c8f65d1406e251df74e72b18913ef48594", // pragma: allowlist secret
                compressedHex: hex"02bfaa0a7de99fcca7702e00abff8bc9be4191bdd3b91d6900c4076c4bfdeb21f6", // pragma: allowlist secret
                evmAddress: 0x82512864E8FECD296D228804F916601e925c97a3
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"04761e5750ce3dd0cdca272572746151200a6216413f3a59681751d04ad4aeb46f78251d54d47dba45624d484d9c002ad007fa62180159bf0dfed6460a5a649765", // pragma: allowlist secret
                compressedHex: hex"03761e5750ce3dd0cdca272572746151200a6216413f3a59681751d04ad4aeb46f", // pragma: allowlist secret
                evmAddress: 0x89B59cDb363602006dA41f2e22DB180c08f08a8c
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"048f9ea6b774a6bb7b776f0d8a10a12d611118de57ed308bb14291668a4815f716976758014910c6e3f518aba17e11d7ee9cca8b91898ae991478557a08791d50f", // pragma: allowlist secret
                compressedHex: hex"038f9ea6b774a6bb7b776f0d8a10a12d611118de57ed308bb14291668a4815f716", // pragma: allowlist secret
                evmAddress: 0xedfe6037703C34cfF2caA52Dc36aD008cF494a8B
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"04c4c09d9163b8a9fe1e64b94e203a632933ed5eecd8bd148c51d665fd75ff2d4afe8b1b2e58b5e5dacfc16b9ad233b9c7c172b35c7925b2ce58e21f9614e7ee06", // pragma: allowlist secret
                compressedHex: hex"02c4c09d9163b8a9fe1e64b94e203a632933ed5eecd8bd148c51d665fd75ff2d4a", // pragma: allowlist secret
                evmAddress: 0xa9C9cb27A16a2adE746fC88D92AE9Bac210d1D87
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"0479c448c6da7a8d5ee2eaa30318d4f5f9e7bde1d804d966d306cbc88e0e85bc0020a6f288651f51abf5caa47e7934de9f0e6dee0888ecb2624c054c4aef5079bf", // pragma: allowlist secret
                compressedHex: hex"0379c448c6da7a8d5ee2eaa30318d4f5f9e7bde1d804d966d306cbc88e0e85bc00", // pragma: allowlist secret
                evmAddress: 0x3AcB94A61A9AC3024e1302b8910f4fA1e647866C
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"0472ec9c47ae07959c59c92b83b03bd1cf393973a4ff9cafedb77e62be405c26d011d660c5caf70e06dbd4f015449ae3518fb99df1c11dc21ebedadc0170c3940f", // pragma: allowlist secret
                compressedHex: hex"0372ec9c47ae07959c59c92b83b03bd1cf393973a4ff9cafedb77e62be405c26d0", // pragma: allowlist secret
                evmAddress: 0x6c18582A638d0DFc4f7da3a4ED9f73FD10EFD5bd
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"04ce0f23517b664591d061a3a3b5af80c3c887abbbbe54f47e3946f3f6405910a62a7b1ee97a6ea45c66aae8263ece2ded5c54ca7a598b6df6b219259a570d4e05", // pragma: allowlist secret
                compressedHex: hex"03ce0f23517b664591d061a3a3b5af80c3c887abbbbe54f47e3946f3f6405910a6", // pragma: allowlist secret
                evmAddress: 0xA6AB5320cC27c4643F946247D2f044946DECF050
            })
        );
    }
}
