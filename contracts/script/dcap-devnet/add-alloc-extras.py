#!/usr/bin/env python3
"""Add CREATE2 deployer and daimo P256Verifier to genesis alloc.

Usage: python3 add-alloc-extras.py <alloc-file>

These contracts are required for DCAP deployment but are not included by
GenerateAlloc.s.sol's vm.etch (which doesn't persist to the alloc dump).
"""

import json
import os
import sys

CREATE2_DEPLOYER = "0x4e59b44847b379578588920cA78FbF26c0B4956C"
CREATE2_CODE = "0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe03601600081602082378035828234f58015156039578182fd5b8082525050506014600cf3"

# daimo P256Verifier — deployed via CREATE2 with zero salt.
# Required by DCAP attestation contracts for ECDSA P256 signature verification.
# Without this, cert uploads and DKG registration fail with empty revert 0x.
P256_VERIFIER = "0xc2b78104907F722DABAc4C69f826a522B2754De4"


def main():
    if len(sys.argv) != 2:
        print(f"Usage: {sys.argv[0]} <alloc-file>")
        sys.exit(1)

    alloc_file = sys.argv[1]
    script_dir = os.path.dirname(os.path.abspath(__file__))

    with open(alloc_file) as f:
        alloc = json.load(f)

    # Load P256Verifier bytecode from file (avoids string truncation issues)
    p256_path = os.path.join(script_dir, "p256verifier_bytecode.txt")
    with open(p256_path) as f:
        p256_code = f.read().strip()

    alloc[CREATE2_DEPLOYER] = {"balance": "0x0", "code": CREATE2_CODE}
    alloc[P256_VERIFIER] = {"balance": "0x0", "code": p256_code}

    with open(alloc_file, "w") as f:
        json.dump(alloc, f, indent=2)

    print(f"Added CREATE2 deployer + P256Verifier ({len(p256_code)} chars). Total: {len(alloc)} accounts")


if __name__ == "__main__":
    main()
