Generate programming problems for my coding platform.

Return ONLY valid JSON. No Markdown fences. No prose before or after JSON.

The top-level result must be an array. Generate exactly the number of problems requested.

Each element MUST have exactly these three properties:

{
"problem": {...},
"tags": [...],
"testCases": [...]
}

## Problem schema

Each `problem` must have exactly:

{
"title": "string",
"slug": "string",
"difficulty": "easy | medium | hard",
"description": "markdown string",
"functionName": "string",
"params": [
{
"name": "string",
"type": "string"
}
],
"returnType": "string"
}

Use camelCase for EVERY JSON property. Never use snake_case.

Required property names:

- functionName
- returnType
- testCases
- expectedOutput
- isSample

Do not use:

- function_name
- return_type
- test_cases
- expected_output
- is_sample

## Problem requirements

`title`:

- Concise problem title.

`slug`:

- Lowercase kebab-case.
- Derived from title.
- Unique across generated problems.

`difficulty`:

- Exactly `easy`, `medium`, or `hard`.

`description`:

- Complete problem statement written in Markdown.
- Explain the task.
- Explain parameters.
- Explain return value.
- Include constraints.
- Include examples when useful.
- Do not describe stdin/stdout.
- Do not include solution code.
- Do not include implementation hints unless they are part of the problem itself.

IMPORTANT:
The description is stored inside JSON, so it MUST contain valid JSON escaping.

Never output invalid JSON escapes such as:

- `\le`
- `\ge`
- `\text`
- `\in`

If using LaTeX, escape every backslash as `\\`.

Prefer plain Markdown operators instead of LaTeX:

- `<=`
- `>=`

For example, use:

`1 <= len(nums) <= 10^5`

instead of:

`$1 \le \text{len}(nums) \le 10^5$`

## Function signature

Problems are solved by calling a function.

`functionName` must be a valid Python function name.

Supported types:

int
float
string
bool
int[]
float[]
string[]
bool[]

Nested arrays are allowed, for example:

int[][]

`returnType` uses the same type syntax.

Keep the function signature consistent with every test case.

Example:

{
"functionName": "twoSum",
"params": [
{
"name": "nums",
"type": "int[]"
},
{
"name": "target",
"type": "int"
}
],
"returnType": "int[]"
}

## Python execution model

The user's solution is inserted into this harness:

from typing import List, Dict, Optional, Tuple, Set
from collections import *

import json
import sys

<USER_CODE>

if **name** == "**main**":
args = json.loads(sys.stdin.read())
result = <FUNCTION_NAME>(*args)
print(json.dumps(result))

Therefore every `testCases[].input` must be a JSON-encoded array containing function arguments in exactly the same order as `params`.

For example:

{
"input": "[[2,7,11,15],9]",
"expectedOutput": "[0,1]"
}

represents:

functionName([2, 7, 11, 15], 9)

`expectedOutput` must be a valid JSON string representing the function's return value.

Examples:

int -> "42"
float -> "3.14"
string -> "\"hello\""
bool -> "true"
int[] -> "[1,2,3]"
string[] -> "[\"a\",\"b\"]"
null -> "null"

## Test cases

Each problem must have at least 8 test cases.

Include:

- 2–3 representative sample cases
- normal cases
- edge cases
- boundary cases
- cases targeting common incorrect solutions

Each test case must have exactly:

{
"input": "JSON string",
"expectedOutput": "JSON string",
"isSample": true,
"ord": 0
}

`isSample` must be `true` for 2–3 representative examples and `false` for all other cases.

`ord` starts at 0 and increments sequentially.

Every test case MUST:

- Match declared parameter types.
- Contain arguments in correct parameter order.
- Contain valid JSON.
- Have correct expected output.
- Respect problem constraints.

## Tags

Each problem must have 2–5 relevant tags.

Each tag must have exactly:

{
"slug": "string",
"name": "string"
}

Use standard algorithm/data-structure concepts where applicable, such as:

arrays
strings
hash-table
sorting
binary-search
two-pointers
sliding-window
stack
queue
tree
graph
dynamic-programming
greedy
recursion
backtracking

## Uniqueness

When generating multiple problems:

- Do not duplicate problems.
- Do not reuse the same title.
- Do not reuse the same slug.
- Avoid generating problems that are essentially the same problem with different wording.
- Keep each problem's function signature, description, and test cases internally consistent.

## Final validation

Before returning the JSON, verify all of the following:

1. Top-level result is an array.
2. Correct number of problems was generated.
3. Every array element contains exactly `problem`, `tags`, and `testCases`.
4. Every JSON property uses camelCase.
5. No snake_case properties exist anywhere.
6. Every problem has exactly the required problem properties.
7. Every function name is valid Python.
8. Every parameter type is supported.
9. `returnType` is supported.
10. Description matches function signature.
11. Every test case matches function signature.
12. Every `input` is valid JSON.
13. Every `expectedOutput` is valid JSON.
14. Every expected output is correct.
15. `ord` values are sequential starting from 0.
16. Each problem has at least 8 test cases.
17. Each problem has 2–3 sample test cases.
18. Constraints match generated test cases.
19. Tags are relevant.
20. No duplicate problems.
21. No invalid JSON escape sequences.
22. Output contains no prose outside the JSON.
23. Output contains no Markdown code fence.

Generate fresh problems unless I explicitly provide a topic, difficulty, algorithm, or other constraint.
