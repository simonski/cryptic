package main

// VERSION is the number of this beast
const VERSION = "0.0.5"

// KP_FILE the key for the env var pointint to the file we load/save
const KP_FILE = "KP_FILE"

// KP_ENCRYPTION_ENABLED the key to the public key
const KP_ENCRYPTION = "KP_ENCRYPTION"

// KP_PUBLIC_KEY the key to the public key
const KP_PUBLIC_KEY = "KP_PUBLIC_KEY"

// KP_PUBLIC_KEY the key to the private key
const KP_PRIVATE_KEY = "KP_PRIVATE_KEY"

// GLOBAL_USAGE - well, it tells me what to type
const GLOBAL_USAGE = `kp is a tool for using key/pairs.

Usage:

    kp <command> [arguments]

The commands are:

    ls                                          list keys
    put <key> (-value VALUE) (-d description)   save "key/value" (read stdin if "-value" is unspecified)
    get <key> (-stdout)                         retrieve key/value to clipboard (or -stdout)
    rm <key>                                    permanently remove "key"

    info                                        review environment variables used
    verify                                      check encryption keys exist and work
    clear                                       remove all values
    version                                     print application version

`

const GLOBAL_SSH_KEYGEN_USAGE = `The following will create a key/pair for encryption: 

     ssh-keygen -m pem -f ~/.ssh/id_rsa
     ssh-keygen -f ~/.ssh/id_rsa.pub -e -m pem > ~/.ssh/id_rsa.pem

You can optionally use environment variables to override the defaults:

     export KP_FILE=~/.kpfile
     export KP_PUBLIC_KEY=~/.ssh/id_rsa.pem
     export KP_PRIVATE_KEY=~/.ssh/id_rsa

`
