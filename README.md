<a href="compiler.go"><img width="731" alt="THE SUPER TINY COMPILER" src="https://cloud.githubusercontent.com/assets/952783/14413766/134c4068-ff39-11e5-996e-9452973299c2.png"/></a>

Inspired by [The Super Tiny Compiler][1] project. I rewrite it by Golang to make it easier to understand, and I also reduce some concept!

### Usage

```
$ git clone git@github.com:ClarkThan/the-super-super-tiny-compiler.git
$ cd the-super-super-tiny-compiler && go build compiler.go

$ ./compiler
add(3, sub(4, len("foo")));

$ ./compiler 23
23;

$ ./compiler '(add 23 (len "Jordan"))'
add(23, len("Jordan"));

$ ./compiler '(foo "Jordan") (add 23 45)'
foo("Jordan");
add(23, 45);
```


### Funny

Feel free to modify code for inspecting the Tokenizer, Parser and CodeGen output.
And I also encourage you improve the error message of the invalid input.

---

[![cc-by-4.0](https://licensebuttons.net/l/by/4.0/80x15.png)](http://creativecommons.org/licenses/by/4.0/)

[1]: https://github.com/thejameskyle/the-super-tiny-compiler