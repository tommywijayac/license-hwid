# Steps
## Setup (one-time)
1. Generate `gen`
```
cd generator && mkdir -m 777 bin && go build -o bin/gen .
```

3. Create your RSA key pairs
```
./bin/gen -rsa
```

## Usage
TODO: should open machine id as args
TODO: should add instruction how to get target machine id

1. Generate license for `machine-a`
```
./bin/gen -generate -hwlabel=machine-a
```

2. Move RSA public key and license to target machine

3. Adjust software code
```
import (
  lic "github.com/tommywijayac/license-hwid"
)

func main() {
  want := []byte(``) // hardcoded expected RSA public key
  isParameterValid, err := lic.ValidatePublicKey(want, "path/to/public-key")
  if err != nil || !isParameterValid {
    // do your thing
    os.Exit(1)
  }

  isLicValid, err := lic.ValidateLicense("path/to/public-key", "path/to/license")
  if err != nil || !isLicValid {
    // do your thing
    os.Exit(1)
  }

  // valid, continue...
}
```
