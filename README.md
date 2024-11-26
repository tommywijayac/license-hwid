# Steps
## Setup (one-time)
1. Generate `gen`
```
go build -o gen .
```

3. Create your RSA key pairs
```
./gen -rsa
```

## Usage
1. Generate license for `machine-a`
```
./gen -generate -hwlabel=machine-a
```

3. Move RSA public key and license to target machine

4. Adjust software code
```
import (
  lic "github.com/tommywijayac/license-hwid"
)

func main() {
  want := []byte(``) // hardcoded expected RSA public key
  isParameterValid, err := lic.ValidatePublicKey(want, "path/to/public-key")
  if !isParameterValid {
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
