// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.23;

/* solhint-disable max-line-length */

contract ValidatorData {
    struct ValidatorInfo {
        bytes uncompressedHex;
        address evmAddress;
    }

    ValidatorInfo[] validators;

    constructor() {
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"04fcb13cd71a77d93593044e8d1f2ea828fb8efc754d92f290443c3a75293ee327f5e1d4ada8b9f4ca123dea7a09c999506819740485b9fb2851ae9d878f54da71", // pragma: allowlist secret
                evmAddress: 0x6C09bCE676733b8130f230682c85b9C0a55DFB1b
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"04bf740fb39d2b61b08929e4821f2a50cda3abef574d93f5387c92ed8ca0dfaadbcaea6a7a6a2759a78976766e42e3219b10d7b03d07aa541deec869107509b5fd", // pragma: allowlist secret
                evmAddress: 0x32ca2A8ED45fD42CF18C7948200E1e17a5bB1BaD
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"04c750cb64f6e4a7ffaaa830decef92b66d5c41785808f1b2580772eb81e16d01dfce9d27bbd85b655fcfb22623c0ff50a98cb90682ecab6907f7db3a5af61db2d", // pragma: allowlist secret
                evmAddress: 0x654c5FCbc74673eB45FdbF7170d8e6b172a21e3d
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"046072c9ef442646eb8c800f4a56a5522ee849a9a9736716b4adbe973187ff450d24e9780461749caf7719d678d56c6945c1cee535a6bb0ee96bd2ffdb918defcc", // pragma: allowlist secret
                evmAddress: 0xb1465CAc452Ec4d718c283Cd9b704049285055A2
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"046e70b9e30395fa9d74135da711229aeeb6aa9e02ea49d0535a2f3b84133706aa80513c3d347a33ab388776d246ca51e742028b4f113c5ddfe8e8e0301c52b76d", // pragma: allowlist secret
                evmAddress: 0x8c1BF40ac76466F987D37D9EcE56379D826C467B
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"042cca88bca3c4a02d0a364bf01aa4068c696efe27a1f372b13dc619c20f021acd3f769411471a42e8592091c030af8bb73075d78cec61b7da2bff05217dc003f8", // pragma: allowlist secret
                evmAddress: 0x56C7b880eD841f031024AFc02a2F96b20479dF45
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"043e175c1cb80c7f562060ab03d1376374f86a595247bdc9ebbcdc479faa41efbe470898e9c86d4ed5fa69a92a02ea00512fd5ab311af01167fb0774a16bc5d098", // pragma: allowlist secret
                evmAddress: 0x14cb3A500bEc93367A685B3eE868f2A0351ea458
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"046d7ee14572ef4a6ea2c55b6e20498f9b6a38ca435de9c649365f8a4f0aba0611942f7770d1ce73a65efd076ebfb7afdb9f6768134c1769ca2f7429abf76baab3", // pragma: allowlist secret
                evmAddress: 0xb6d13E31E121ADA97E825aCe11cf50290Ef4E9d9
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"04cea9d4416bfa1aa8560e83fef4313dd1d74a7c84dff844be2b4af560b446884b987549b74389bb851a28fe2d4ded5200a5da89aeee1807f8646835779079f016", // pragma: allowlist secret
                evmAddress: 0x9753a62659f8F91F0bad8331993C530978BE0221
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"04ee4b72d1aeea63101e64cc4b1e1fa83dd73ddbbf293a6aa544d01c3c569b1a970189c3e7e21644e78c8ef1fb602ca9cf6546d6f2e3f5adf1cba6c663e5854898", // pragma: allowlist secret
                evmAddress: 0xeec378BeAE63bb37d74B0DB85E7A0b52cc4b565e
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"04642dd9f86f9088a7a51b4a0b6c34112ad4454255ec7c32b60cd9c6e0a4cc97abd4c458b512773f581a7d082cd20bcec884520b72c41f46302a757838d47053bd", // pragma: allowlist secret
                evmAddress: 0x67bEC6E10494F1afdBd29461E8E6281E4e9e27a1
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"0427b6000e74b541b97079d94602f40b67631787c647b058750789a90a045a12d40cc32c966e2b94065ce0eb25cbb34ba9f6dfe2a25ed12d175a51fca4c5ad4e9e", // pragma: allowlist secret
                evmAddress: 0x9b02B30662f1F1BD07ED5F92d3a54a9f5cf56feC
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"043e185b04210502055959abac4866f362350e73a44b607c2f2a91d0ef6b2bdb62b164bf81eebcfa9a046035b2cee4894b1711248f1114818c36f14279d4d6c7e0", // pragma: allowlist secret
                evmAddress: 0x950a52f2bebabbd8C5F18BbfD3aAF7E0B6338C6A
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"0434d8a7557e31ef2db308be34b6537652d57a83bf517daf7421e19e0df6cc45e23c4d5816045c4012aceeed1e8c98f1a5134e91bacbefa3cb5196e979a065f349", // pragma: allowlist secret
                evmAddress: 0xa9E6eB0bE1c7af03cF1bB80Ba583bE39927D9E7b
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"0482782124bc966f7fecd03c38aa4cac234dc4e4e3cecf04d579143714e525fa57f7af9475fbaf7fa78ffb975f6d58e245bea952dd039f0fec4e9db418c3b2955d", // pragma: allowlist secret
                evmAddress: 0x5B7CF83241848f5E25155238fEa8FbADdc24038d
            })
        );
        validators.push(
            ValidatorInfo({
                uncompressedHex: hex"049aec970a392847bf70305dfd8c1110d822328e04ed48386d792c9b9d87f8ca6d724ac46bc58ce02e50c03fba18012a7c07b38bc92d3940c088fbb53164689634", // pragma: allowlist secret
                evmAddress: 0xfC75911DBb20508f70641274FbcB334D395AE373
            })
        );
    }
}
