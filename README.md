[![progress-banner](https://backend.codecrafters.io/progress/interpreter/39e99a1a-eff8-41fd-b8e8-4198bd7df651)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)

# Lox Interpreter in Go

This project is a Go implementation of the Lox interpreter, built as part of the [Codecrafters "Build your own Interpreter" challenge](https://app.codecrafters.io/courses/interpreter/overview). The challenge follows the fantastic book [Crafting Interpreters](https://craftinginterpreters.com/) by Robert Nystrom.

## 📌 What This Project Does

Lox is a simple, dynamically-typed scripting language. This interpreter project covers the fundamental stages of language implementation:

- **Scanning (Lexing/Tokenization):** Converting raw source code into a stream of tokens.
- **Parsing:** Building an Abstract Syntax Tree (AST) from the token stream, representing the code\'s structure.
- **Resolution:** Analyzing variable scopes and ensuring correct bindings.
- **Interpretation:** Walking the AST to execute the Lox code.

## ✨ Key Features

- **Scanner (`internal/scanner`):** Transforms Lox source code into a sequence of tokens (e.g., keywords, identifiers, literals, operators).
- **Parser (`internal/parser`):** Constructs an AST from the tokens, defining the grammatical structure of expressions and statements.
- **AST (`internal/ast`):** Defines the `Expr` (expression) and `Stmt` (statement) interfaces and their concrete types that form the nodes of the tree.
- **Resolver (`internal/resolver`):** Performs static analysis to resolve variable lookups, determining how many levels up the environment chain a variable is defined.
- **Interpreter (`internal/interpreter`):** Traverses the AST to execute Lox programs. It manages:
    - Variable declaration, assignment, and lookup.
    - Control flow structures (if/else, while loops).
    - Function definition, calls, and return values.
    - Logical and binary operations.
    - Class declaration, instantiation, and method calls.
    - Object-oriented features including inheritance and the `super` keyword.
- **Error Handling:** Provides descriptive error messages for syntax errors (compile-time) and runtime errors.
- **Environment & Scope (`internal/interpreter/environment.go`):** Manages lexical scoping for variables and functions using a chain of environment objects.
- **Callable Types (`internal/interpreter/callable.go`):** Supports Lox functions and native Go functions callable from Lox.
- **Object-Oriented Programming:** Full support for classes, inheritance, instance methods, constructors, and the `this` and `super` keywords.
- **Native Functions:** Includes a `clock()` native function to get the current system time.

## 🏗️ Object-Oriented Features

The interpreter fully supports object-oriented programming with the following features:

### Classes and Instances
- **Class Declaration:** Define classes with the `class` keyword.
- **Instantiation:** Create instances by calling the class like a function.
- **Instance Properties:** Get and set properties on instances using dot notation.
- **Methods:** Define methods within classes that have access to `this`.

### Inheritance
- **Single Inheritance:** Classes can inherit from a single superclass using the `<` operator.
- **Method Overriding:** Subclasses can override methods from their superclass.
- **Super Keyword:** Access superclass methods using `super.methodName()`.

### Special Methods
- **Constructors:** Define an `init()` method to initialize new instances.
- **This Binding:** The `this` keyword refers to the current instance within methods.

### Example:
```lox
class Vehicle {
  init(brand) {
    this.brand = brand;
  }

  describe() {
    print "This is a " + this.brand + " vehicle";
  }
}

class Car < Vehicle {
  init(brand, model) {
    super.init(brand);
    this.model = model;
  }

  describe() {
    print "This is a " + this.brand + " " + this.model;
  }

  honk() {
    print "Beep beep!";
  }
}

var myCar = Car("Toyota", "Camry");
myCar.describe(); // "This is a Toyota Camry"
myCar.honk();     // "Beep beep!"
```

## 🛠️ Why I Built This Project

Having previously completed this challenge in Java, I wanted to deepen my understanding of interpreters and explore the nuances of implementing one in Go. This project served as a great opportunity to:
- Solidify concepts from "Crafting Interpreters."
- Gain practical experience with Go\'s syntax, and standard library.
- Compare and contrast language implementation approaches between Java and Go.
- Tackle the interesting problems that arise when building a language from scratch.

## 🔍 How It Works Internally

The interpreter processes Lox code in several stages:

1.  **Reading Source:** The `main` function in `cmd/myinterpreter/main.go` reads the Lox source file specified by the user.
2.  **Scanning:** The `scanner.NewScanner(source).ScanTokens()` method breaks the raw source string into a list of `Token` objects. Each token has a type (e.g., `IDENTIFIER`, `STRING`, `NUMBER`, `PLUS`), a lexeme (the actual text), a literal value (for strings and numbers), and a line number for error reporting.
3.  **Parsing:** The `parser.NewParser(tokens).Parse()` method (or `GetStatements()` for a full program) takes the list of tokens and constructs an AST. It uses a recursive descent approach, with separate functions for parsing different grammatical rules (e.g., expressions, statements, primary expressions). If syntax errors are encountered, they are reported.
4.  **Resolving (Static Analysis):** Before interpretation, the `resolver.NewResolver(interpreter).Resolve(statements)` pass walks the AST. Its primary job is to determine, for each variable access, how many environments "up" the chain it needs to go to find the variable\'s definition. This information is stored and used by the interpreter for efficient lookups. It also catches errors like reading a variable in its own initializer or using `return` outside a function.
5.  **Interpreting:** The `interpreter.NewInterpreter().Run(statements)` (or `Interpret()` for single expressions) method walks the AST using the Visitor pattern. Each node type in the AST (e.g., `BinaryExpr`, `IfStmt`, `FunctionStmt`, `ClassStmt`) has a corresponding `Visit...` method in the interpreter.
    *   **Expressions** are evaluated to produce a value.
    *   **Statements** are executed for their side effects (e.g., defining a variable, printing output).
    *   **Environments** are created for blocks and function calls to manage lexical scope. When a variable is looked up, the resolver\'s previously computed depth is used to quickly find it in the correct environment.
    *   **Function calls** involve creating a new environment, binding arguments to parameters, and executing the function\'s body.
    *   **Class instantiation** creates new instances with their own property storage and method binding.
    *   **Inheritance** is handled by copying methods from superclasses and managing the `super` keyword for method calls up the inheritance chain.

## ⚙️ How to Set Up and Run

### Prerequisites
-   Go (version 1.22 or later recommended, as per the original `README.md`). You can check your version with `go version`.

### Installation
1.  Clone the repository:
    ```sh
    git clone https://github.com/Md-Talim/codecrafters-interpreter-go.git
    cd codecrafters-interpreter-go
    ```

### Running the Interpreter

The main entry point is `cmd/myinterpreter/main.go`. A helper script `your_program.sh` is provided to simplify running different modes of the interpreter.

**General Usage:**
```sh
./your_program.sh <command> <filename>
```

**Available Commands:**

| Command    | Description                                                                       | Example                                     |
| :--------- | :-------------------------------------------------------------------------------- | :------------------------------------------ |
| `tokenize` | Scans the source file and prints the stream of tokens.                            | `./your_program.sh tokenize test.lox`       |
| `parse`    | Scans and parses the source file, then prints a string representation of the AST. | `./your_program.sh parse test.lox`          |
| `evaluate` | Scans, parses, and evaluates a *single expression* from the source file.          | `./your_program.sh evaluate expression.lox` |
| `run`      | Scans, parses, resolves, and interprets the entire Lox program in the file.       | `./your_program.sh run test.lox`            |

*(Note: The `evaluate` command might be more suited for a REPL or specific testing scenarios where a file contains only a single expression. For full programs, use `run`.)*

### Example

1.  Create a Lox file, for example, `hello.lox`:
    ```lox
    // hello.lox
    print "Hello, Lox!";

    fun greet(name) {
      print "Hello, " + name + "!";
    }

    greet("User");

    var a = 10;
    if (a > 5) {
      print "a is greater than 5";
    } else {
      print "a is not greater than 5";
    }

    // Class and inheritance example
    class Animal {
      init(name) {
        this.name = name;
      }

      speak() {
        print this.name + " makes a sound";
      }
    }

    class Dog < Animal {
      init(name, breed) {
        super.init(name);
        this.breed = breed;
      }

      speak() {
        print this.name + " barks";
      }

      getBreed() {
        return this.breed;
      }
    }

    var myDog = Dog("Rex", "Golden Retriever");
    myDog.speak();
    print "Breed: " + myDog.getBreed();

    print clock(); // Native function
    ```

2.  Run the interpreter:
    ```sh
    ./your_program.sh run hello.lox
    ```

    Expected Output:
    ```
    Hello, Lox!
    Hello, User!
    a is greater than 5
    Rex barks
    Breed: Golden Retriever
    <current_timestamp>
    ```

## 📁 Folder and File Structure

-   `cmd/myinterpreter/main.go`: The command-line interface for the interpreter.
-   `internal/`: Contains the core logic of the interpreter.
    -   `ast/`: Abstract Syntax Tree definitions (expressions, statements, tokens).
    -   `scanner/`: The lexer/tokenizer.
    -   `parser/`: The parser that builds the AST.
    -   `resolver/`: The variable resolver for static scope analysis.
    -   `interpreter/`: The tree-walk interpreter that executes the AST.
-   `pkg/lox/`: High-level functions that orchestrate the different stages (tokenize, parse, run).
-   `your_program.sh`: Shell script to run the interpreter.
-   `test.lox`: A sample Lox file for testing.
-   `testdata/`: Contains more specific Lox test files and their expected outputs for various language features.

## 💡 Challenges & Lessons Learned

-   **Pointer vs. Value Semantics in Go:** Deciding when to use pointers for AST nodes and other structures was a key consideration, especially for shared data like environments or when modifying state during tree walks.
-   **Error Handling in Go:** Go\'s explicit error handling (returning `error` values) is robust but requires careful propagation. I learned to define custom error types (e.g., `RuntimeError`, `returnValue`) to distinguish different error conditions and control flow.
-   **Visitor Pattern Implementation:** Implementing the Visitor pattern in Go for the AST (for printing, resolving, and interpreting) was a good exercise. Go\'s interfaces make this quite clean.
-   **Managing State:** Keeping track of the current environment, function type (for `return` statement validation), and scope resolution details required careful state management within the Resolver and Interpreter.
-   **Debugging the Parser:** Recursive descent parsers can be tricky to debug. Print statements and a methodical approach were essential. The `AstPrinter` was invaluable.
-   **Testing:** Writing small Lox test files for each feature and comparing the output against expected results was crucial for ensuring correctness.
-   **Class and Inheritance Implementation:** Implementing object-oriented features required careful consideration of method binding, `this` resolution, and the `super` keyword. Managing the inheritance chain and method lookup was particularly challenging.
-   **Memory Management:** With classes and instances, ensuring proper memory management and avoiding circular references became more important, especially when dealing with method closures and `this` binding.


## 🙏 Acknowledgments

-   This project is an implementation of the interpreter described in the book [Crafting Interpreters](https://craftinginterpreters.com/) by Robert Nystrom.
-   The [Codecrafters "Build your own Interpreter" challenge](https://app.codecrafters.io/courses/interpreter/overview) provided the structure and motivation for this project.