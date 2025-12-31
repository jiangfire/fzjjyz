package i18n

// enUS 英文翻译字典.
type enUS struct{}

func (e *enUS) Get(key string) string {
	translation, exists := enTranslations[key]
	if !exists {
		return ""
	}
	return translation
}

// enTranslations 英文翻译映射.
var enTranslations = map[string]string{
	// 根命令和应用信息
	"app.name":        "fzjjyz",
	"app.description": "Post-quantum file encryption tool - Using Kyber768 + ECDH + AES-256-GCM + Dilithium3",
	"app.long": `fzjjyz - Post-quantum file encryption tool

Provides secure file encryption using the following algorithms:
  • Kyber768 - Post-quantum key encapsulation
  • X25519 ECDH - Traditional key exchange
  • AES-256-GCM - Authenticated encryption
  • Dilithium3 - Digital signature

Quick start:
  # Generate key pair
  fzjjyz keygen -d ./keys -n mykey

  # Encrypt file
  fzjjyz encrypt -i plaintext.txt -o encrypted.fzj -p keys/mykey_public.pem -s keys/mykey_dilithium_private.pem

  # Decrypt file
  fzjjyz decrypt -i encrypted.fzj -o decrypted.txt -p keys/mykey_private.pem -s keys/mykey_dilithium_public.pem

  # View file info
  fzjjyz info -i encrypted.fzj

Project home: https://codeberg.org/jiangfire/fzjjyz`,

	// 全局标志
	"flags.verbose": "Enable verbose output",
	"flags.force":   "Force overwrite existing files",

	// encrypt 命令
	"encrypt.short": "Encrypt file",
	"encrypt.long": `Encrypt file using post-quantum hybrid encryption algorithm.

Encryption process:
  1. Read original file
  2. Generate random session key
  3. Kyber768 + ECDH key encapsulation
  4. AES-256-GCM data encryption
  5. Dilithium3 signature verification
  6. Build encrypted file header
  7. Write encrypted file

Required parameters:
  --input, -i         Input file path
  --public-key, -p    Kyber+ECDH public key file
  --sign-key, -s      Dilithium private key file

Examples:
  fzjjyz encrypt -i plaintext.txt -o encrypted.fzj -p public.pem -s dilithium_private.pem
  fzjjyz encrypt --input data.txt --public-key pub.pem --sign-key priv.pem --force`,
	"encrypt.flags.input":       "Input file path (required)",
	"encrypt.flags.output":      "Output file path (optional, default: input.fzj)",
	"encrypt.flags.public-key":  "Kyber+ECDH public key file (required)",
	"encrypt.flags.sign-key":    "Dilithium private key file (required)",
	"encrypt.flags.force":       "Overwrite output file",
	"encrypt.flags.buffer-size": "Buffer size (KB), 0=auto",
	"encrypt.flags.streaming":   "Use streaming mode (recommended for large files)",

	// decrypt 命令
	"decrypt.short": "Decrypt file",
	"decrypt.long": `Decrypt file using post-quantum hybrid encryption algorithm.

Decryption process:
  1. Parse file header
  2. Verify file format
  3. Kyber768 + ECDH key decapsulation
  4. AES-256-GCM data decryption
  5. Verify SHA256 hash
  6. Verify Dilithium signature
  7. Write original file

Required parameters:
  --input, -i         Encrypted file path
  --private-key, -p   Kyber+ECDH private key file
  --verify-key, -s    Dilithium public key file (optional)

Examples:
  fzjjyz decrypt -i encrypted.fzj -o decrypted.txt -p private.pem -s dilithium_public.pem
  fzjjyz decrypt --input data.fzj --private-key priv.pem --verify-key pub.pem --force`,
	"decrypt.flags.input":       "Encrypted file path (required)",
	"decrypt.flags.output":      "Output file path (optional, default: original filename)",
	"decrypt.flags.private-key": "Kyber+ECDH private key file (required)",
	"decrypt.flags.verify-key":  "Dilithium public key file (optional)",
	"decrypt.flags.force":       "Overwrite output file",
	"decrypt.flags.buffer-size": "Buffer size (KB), 0=auto",
	"decrypt.flags.streaming":   "Use streaming mode (recommended for large files)",

	// encrypt-dir 命令
	"encrypt-dir.short": "Encrypt directory",
	"encrypt-dir.long": `Pack entire directory into ZIP, then encrypt using post-quantum hybrid encryption.

Encryption process:
  1. Scan source directory recursively
  2. Pack directory structure into ZIP format
  3. Read ZIP data into memory
  4. Kyber768 + ECDH key encapsulation
  5. AES-256-GCM ZIP data encryption
  6. Dilithium3 signature verification
  7. Build encrypted file header
  8. Write encrypted file (.fzj)

Required parameters:
  --input, -i         Source directory path
  --output, -o        Output encrypted file path
  --public-key, -p    Kyber+ECDH public key file
  --sign-key, -s      Dilithium private key file

Examples:
  fzjjyz encrypt-dir -i ./sensitive_data -o secure.fzj -p public.pem -s dilithium_private.pem
  fzjjyz encrypt-dir --input ./confidential --output backup.fzj --public-key pub.pem --sign-key priv.pem --force`,
	"encrypt-dir.flags.input":       "Source directory path (required)",
	"encrypt-dir.flags.output":      "Output encrypted file path (required)",
	"encrypt-dir.flags.public-key":  "Kyber+ECDH public key file (required)",
	"encrypt-dir.flags.sign-key":    "Dilithium private key file (required)",
	"encrypt-dir.flags.force":       "Overwrite output file",
	"encrypt-dir.flags.buffer-size": "Buffer size (KB), 0=auto",
	"encrypt-dir.flags.streaming":   "Use streaming mode",

	// decrypt-dir 命令
	"decrypt-dir.short": "Decrypt directory",
	"decrypt-dir.long": `Decrypt encrypted directory archive and restore original directory structure.

Decryption process:
  1. Parse encrypted file header
  2. Verify file format
  3. Kyber768 + ECDH key decapsulation
  4. AES-256-GCM data decryption
  5. Verify SHA256 hash
  6. Verify Dilithium signature
  7. Extract ZIP to target directory
  8. Restore original directory structure

Required parameters:
  --input, -i         Encrypted file path (.fzj)
  --output, -o        Output directory path
  --private-key, -p   Kyber+ECDH private key file
  --verify-key, -s    Dilithium public key file (optional)

Examples:
  fzjjyz decrypt-dir -i secure.fzj -o ./restored -p private.pem -s dilithium_public.pem
  fzjjyz decrypt-dir --input backup.fzj --output ./recovered --private-key priv.pem --verify-key pub.pem --force`,
	"decrypt-dir.flags.input":       "Encrypted file path (required)",
	"decrypt-dir.flags.output":      "Output directory path (required)",
	"decrypt-dir.flags.private-key": "Kyber+ECDH private key file (required)",
	"decrypt-dir.flags.verify-key":  "Dilithium public key file (optional)",
	"decrypt-dir.flags.force":       "Force overwrite existing files in output directory",
	"decrypt-dir.flags.buffer-size": "Buffer size (KB), 0=auto",
	"decrypt-dir.flags.streaming":   "Use streaming mode",

	// keygen 命令
	"keygen.short": "Generate post-quantum key pair",
	"keygen.long": `Generate complete key pair combination:
  • Kyber768 + ECDH key pair (for encryption/decryption)
  • Dilithium3 key pair (for signing/verification)

Generated files:
  {name}_public.pem          - Kyber+ECDH public key
  {name}_private.pem         - Kyber+ECDH private key (0600 permissions)
  {name}_dilithium_public.pem  - Dilithium public key
  {name}_dilithium_private.pem - Dilithium private key (0600 permissions)

Examples:
  fzjjyz keygen -d ./keys -n mykey
  fzjjyz keygen --output-dir ./keys --name mykey --force`,
	"keygen.flags.output-dir": "Output directory",
	"keygen.flags.name":       "Key name prefix (default: timestamp)",
	"keygen.flags.force":      "Overwrite existing files",

	// keymanage 命令
	"keymanage.short": "Key management tool",
	"keymanage.long": `Manage encryption keys, supporting export, import, and verification operations.

Available operations:
  export    Extract and export public key from private key file
  import    Import key files to specified directory
  verify    Verify key pair matching

Examples:
  # Export public key
  fzjjyz keymanage export --private-key private.pem --output public_extracted.pem

  # Verify key pair
  fzjjyz keymanage verify --public-key public.pem --private-key private.pem

  # Import keys
  fzjjyz keymanage import --public-key pub.pem --private-key priv.pem --output-dir ./keys`,
	"keymanage.flags.action":      "Action type: export/import/verify (required)",
	"keymanage.flags.public-key":  "Public key file path",
	"keymanage.flags.private-key": "Private key file path",
	"keymanage.flags.output":      "Output file path (for export)",
	"keymanage.flags.output-dir":  "Output directory (for import)",

	// info 命令
	"info.short": "View encrypted file information",
	"info.long": `Parse and display detailed information about encrypted files, including:
  • Filename and original size
  • Encryption timestamp
  • Algorithms used
  • Signature status
  • Integrity verification

Examples:
  fzjjyz info -i encrypted.fzj
  fzjjyz info --input data.fzj`,
	"info.flags.input": "Encrypted file path (required)",

	// version 命令
	"version.short":       "Show version information",
	"version.long":        "Show fzjjyz version and build details",
	"version.info":        "Version information",
	"version.label":       "Version:",
	"version.app_name":    "App name:",
	"version.description": "Description:",

	// Progress messages
	"progress.loading_keys":         "Loading keys...",
	"progress.encrypting":           "Encrypting file...",
	"progress.verifying":            "Verifying...",
	"progress.decrypting":           "Decrypting file...",
	"progress.packing":              "Packing directory...",
	"progress.extracting":           "Extracting files...",
	"progress.generating_kyber":     "Generating Kyber768 keys...",
	"progress.generating_ecdh":      "Generating ECDH X25519 keys...",
	"progress.generating_dilithium": "Generating Dilithium3 signature keys...",
	"progress.saving_keys":          "Saving key files...",

	// Status messages
	"status.done":   "Done",
	"status.failed": "Failed",
	"status.warning_no_sign_verify": "⚠️  Warning: No signature verification key " +
		"provided, skipping signature verification",
	"status.success_encrypt": "✅ Encryption successful!",
	"status.success_decrypt": "✅ Decryption successful!",
	"status.success_keygen":  "✅ Key pair generated successfully!",
	"status.success_export":  "✅ Public key exported to: %s",
	"status.success_import":  "✅ Keys imported to: %s",
	"status.success_verify":  "✅ Key pair verified",
	"status.failed_verify":   "❌ Key pair mismatch",
	"status.encrypting_file": "Encrypting file: %s",
	"status.decrypting_file": "Decrypting file: %s",
	"status.encrypting_dir":  "Encrypting directory: %s",
	"status.decrypting_dir":  "Decrypting directory: %s",
	"status.generating_keys": "Generating key pair...",
	"status.public_key":      "Public key",
	"status.sign_key":        "Sign key",
	"status.streaming_mode":  "Streaming mode",

	// File info output
	"file_info.header":            "📁 File info: %s",
	"file_info.basic":             "Basic information:",
	"file_info.encryption":        "Encryption information:",
	"file_info.keys":              "Key information:",
	"file_info.integrity":         "Integrity:",
	"file_info.verification":      "Verification status:",
	"file_info.original_file":     "Original file: %s (%d bytes)",
	"file_info.encrypted_file":    "Encrypted file: %s (%d bytes)",
	"file_info.decrypted_file":    "Decrypted file: %s (%d bytes)",
	"file_info.compressed_rate":   "Compression rate: %.1f%%",
	"file_info.timestamp":         "Timestamp: %s",
	"file_info.algorithm":         "Algorithm: %s (0x%02x)",
	"file_info.version":           "Version: 0x%04x",
	"file_info.magic":             "Magic: %c%c%c\\x%02x",
	"file_info.kyber":             "Kyber encapsulation: %d bytes",
	"file_info.ecdh":              "ECDH public key: %d bytes",
	"file_info.iv":                "IV/Nonce: %d bytes",
	"file_info.signature":         "Signature: %d bytes",
	"file_info.hash":              "SHA256 hash: %x...",
	"file_info.signature_status":  "Signature:",
	"file_info.data_integrity":    "Data integrity:",
	"file_info.exists":            "Present",
	"file_info.not_exists":        "Absent",
	"file_info.complete":          "Complete",
	"file_info.incomplete":        "Incomplete",
	"file_info.original_filename": "Original filename: %s",
	"file_info.file_count":        "File count: %d",
	"file_info.source_dir":        "Source directory: %s",
	"file_info.output_dir":        "Output directory: %s",
	"file_info.zip_size":          "ZIP size: %d bytes",
	"file_info.decrypted_size":    "Decrypted size: %d bytes",
	"file_info.buffer_size":       "Buffer size: %d KB",

	// Directory encryption/decryption info
	"dir_info.encrypt_summary": `File information:
  Source directory: %s
  File count: %d
  ZIP size: %d bytes
  Encrypted file: %s (%d bytes)
  Compression rate: %.1f%%`,
	"dir_info.decrypt_summary": `File information:
  Encrypted file: %s (%d bytes)
  Decrypted size: %d bytes
  File count: %d
  Output directory: %s
  Original filename: %s
  Timestamp: %s`,

	// Single file encryption/decryption info
	"file_info.encrypt_summary": `File information:
  Original file: %s (%d bytes)
  Encrypted file: %s (%d bytes)
  Compression rate: %.1f%%`,
	"file_info.decrypt_summary": `File information:
  Encrypted file: %s (%d bytes)
  Decrypted file: %s (%d bytes)
  Original filename: %s
  Timestamp: %s`,

	// Key generation info
	"keygen_info.files": `Generated files:
  • %s (public key)
  • %s (private key - 0600 permissions)
  • %s (signature public key)
  • %s (signature private key - 0600 permissions)`,

	// Key verification info
	"keymanage_verify.kyber": "  Kyber:  %s",
	"keymanage_verify.ecdh":  "  ECDH:   %s",

	// Security warnings
	"security.warning":        "⚠️  Security warning:",
	"security.protect_keys":   "• Please keep private key files secure",
	"security.no_sharing":     "• Do not share private keys with others",
	"security.secure_storage": "• Use secure storage media",

	// Archive info
	"archive.packed":    "Done (size: %d bytes, files: %d)",
	"archive.decrypted": "Done (size: %d bytes)",

	// Error messages - File related
	"error.file_not_exists":           "File not found: %s",
	"error.input_file_not_exists":     "Input file not found: %s",
	"error.encrypted_file_not_exists": "Encrypted file not found: %s",
	"error.source_dir_not_exists":     "Source directory not found: %s",
	"error.input_not_dir":             "Input path is not a directory: %s",
	"error.output_not_dir":            "Output path is not a directory: %s",
	"error.output_file_exists":        "Output file already exists: %s (use --force to overwrite)",
	"error.output_dir_not_empty":      "Output directory not empty: %s (use --force to overwrite)",
	"error.cannot_create_dir":         "Cannot create directory %s: %v",
	"error.cannot_open_file":          "Cannot open encrypted file: %v",
	"error.cannot_read_file":          "Cannot read file: %v",
	"error.cannot_read_data":          "Cannot read decrypted data: %v",
	"error.cannot_open_temp":          "Cannot open temporary file: %v",

	// Error messages - Key related
	"error.key_not_found": "Key file not found: %s",
	"error.key_invalid":   "Invalid key format",
	"error.load_public_key_failed": `❌ Failed to load public key: %v

Tips:
  1. Check public key file path: %s
  2. Ensure key format is correct (PEM format)
  3. Check file permissions (must be readable)
  4. If first use, generate key pair first: fzjjyz keygen`,
	"error.load_private_key_failed": `❌ Failed to load private key: %v

Tips:
  1. Check private key file path: %s
  2. Ensure key format is correct (PEM format)
  3. Check file permissions (recommended 0600)
  4. Private key should be readable only by owner
  5. Ensure using matching key from encryption`,
	"error.load_sign_key_failed": `❌ Failed to load signature private key: %v

Tips:
  1. Check Dilithium private key file path: %s
  2. Ensure key format is correct (PEM format)
  3. Check file permissions (recommended 0600)
  4. Private key should be readable only by owner
  5. If first use, generate key pair first: fzjjyz keygen`,
	"error.load_verify_key_failed": `❌ Failed to load verification public key: %v

Tips:
  1. Check Dilithium public key file path: %s
  2. Ensure key format is correct (PEM format)
  3. Check file permissions (must be readable)
  4. Ensure using matching key from encryption
  5. If not provided, can omit this parameter (but signature won't be verified)`,
	"error.keygen_kyber_failed":     "Kyber key generation failed: %v",
	"error.keygen_ecdh_failed":      "ECDH key generation failed: %v",
	"error.keygen_dilithium_failed": "Dilithium key generation failed: %v",
	"error.save_keys_failed":        "Failed to save key files: %v",
	"error.save_dilithium_failed":   "Failed to save Dilithium keys: %v",
	"error.export_key_failed":       "Failed to export public key: %v",
	"error.save_export_failed":      "Failed to save public key file: %v",
	"error.import_keys_failed":      "Failed to import keys: %v",
	"error.verify_keys_failed":      "Key pair mismatch",

	// Error messages - Encryption/Decryption related
	"error.encrypt_failed": `❌ Encryption failed: %v

Possible causes:
  1. Insufficient file permissions (cannot read input or write output)
  2. Insufficient memory (large files require more memory)
  3. Key mismatch
  4. Input file modified during encryption

Suggestions:
  - Check disk space and file permissions
  - For large files, try adjusting buffer size with --buffer-size
  - Ensure keys match correctly`,
	"error.decrypt_failed": `❌ Decryption failed: %v

Possible causes:
  1. Key mismatch (using wrong private key)
  2. File corrupted or tampered
  3. Incorrect file format (not fzjjyz encrypted)
  4. Signature verification failed (file may be tampered)
  5. Insufficient file permissions

Security tips:
  - If hash mismatch, file may be corrupted, do not use
  - If signature invalid, keys may not match or file modified
  - Always provide signature verification key for integrity`,
	"error.pack_failed": `❌ Packing failed: %v

Possible causes:
  1. Insufficient directory permissions
  2. Contains unsupported file types (e.g., symlinks)
  3. Insufficient disk space`,
	"error.extract_failed": `❌ Extraction failed: %v

Possible causes:
  1. Insufficient output directory permissions
  2. Insufficient disk space
  3. ZIP file corrupted`,
	"error.temp_file_failed":       "❌ Failed to create temporary file: %v",
	"error.parse_header_failed":    "Failed to parse file header: %v",
	"error.validate_header_failed": "Failed to validate file header: %v",

	// Error messages - Other
	"error.unknown_action":         "Unknown action: %s (supported: export, import, verify)",
	"error.missing_required_flags": "Must provide %s",
	"error.missing_both_keys":      "Must provide --public-key and --private-key",
	"error.nothing_to_do":          "Nothing to do",
}
