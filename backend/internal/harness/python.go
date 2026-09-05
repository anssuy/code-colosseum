package harness

import "fmt"

const pythonPrelude = "from typing import List, Dict, Optional, Tuple, Set\nfrom collections import *\n\n"

func Python(sig Signature, userCode string) string {
	return fmt.Sprintf(`%simport json
import sys

%s

if __name__ == "__main__":
    args = json.loads(sys.stdin.read())
    result = %s(*args)
    print(json.dumps(result))
`, pythonPrelude, userCode, sig.FunctionName)
}
