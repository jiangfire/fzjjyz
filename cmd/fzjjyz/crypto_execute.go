package main

import (
	"codeberg.org/jiangfire/fzjjyz/cmd/fzjjyz/utils"
	"codeberg.org/jiangfire/fzjjyz/internal/zjcrypto"
	"github.com/cloudflare/circl/sign/dilithium/mode3"
)

func calculateBufferSizeFromFile(path string, overrideKB int) int {
	if overrideKB > 0 {
		return overrideKB * 1024
	}
	size, err := utils.GetFileSize(path)
	if err != nil {
		return zjcrypto.OptimalBufferSize(0)
	}
	return zjcrypto.OptimalBufferSize(size)
}

func calculateBufferSizeFromLength(size int64, overrideKB int) int {
	if overrideKB > 0 {
		return overrideKB * 1024
	}
	return zjcrypto.OptimalBufferSize(size)
}

func runEncryptWithMode(
	inputPath, outputPath string,
	hybridPub *zjcrypto.HybridPublicKey,
	dilithiumPriv *mode3.PrivateKey,
	streaming bool,
	bufferSize int,
) error {
	if streaming {
		return zjcrypto.EncryptFileStreaming(
			inputPath, outputPath,
			hybridPub.Kyber, hybridPub.ECDH,
			dilithiumPriv,
			bufferSize,
		)
	}
	return zjcrypto.EncryptFile(
		inputPath, outputPath,
		hybridPub.Kyber, hybridPub.ECDH,
		dilithiumPriv,
	)
}

func runDecryptWithMode(
	inputPath, outputPath string,
	hybridPriv *zjcrypto.HybridPrivateKey,
	dilithiumPub *mode3.PublicKey,
	streaming bool,
	bufferSize int,
) error {
	if streaming {
		return zjcrypto.DecryptFileStreaming(
			inputPath, outputPath,
			hybridPriv.Kyber, hybridPriv.ECDH,
			dilithiumPub,
			bufferSize,
		)
	}
	return zjcrypto.DecryptFile(
		inputPath, outputPath,
		hybridPriv.Kyber, hybridPriv.ECDH,
		dilithiumPub,
	)
}
