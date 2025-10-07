> [!WARNING]
> The library is not in contributions ready state yet , test not exist

# HtmlPuzzles

**HtmlPuzzles** is an experimental Go library for creating custom HTML components with a custom render pipeline. This project was created for educational purposes to study HTML parsing, DSL creation, and architectural patterns.

##  Main Idea

The library allows you to:
- Create custom HTML tags with custom logic
- Implement your own components (loops, conditions, variables, etc.)
- Manage execution context with layer support
- Parse and render HTML with custom tags

##  Installation

```bash
go get github.com/DilemaFixer/HtmlPuzzles
```

For more details chack [wiki page](https://github.com/DilemaFixer/HtmlPuzzles/wiki)

## Example of possibilities

Input html : 
``` Html
<!-- <sync> - start rendering inner html in separate gorutine without stoping main -->
<!-- <for> - just repiting inner html , requared 'itr_count' attrebute  -->
<!-- <templ> - render .html file as innner content -->
<!-- <wrapper> - using as templ for html tag -->
<!-- 'wrapped' attrebute requared always -->
<!-- all attrebutes will be move to tag that describe in wrapped -->
<!-- if attrebute start with ':" prefix will be geting from runtime context , you can describe way to prop -->
<!-- :src="user.AvatarUrl|string" -> src="path/to/img/from/prop" -->
<div>
    <sync>
        <for itr_count=3>
            <h1>Hello World!</h1>
        </for>
    </sync>
    <sync>
        <templ source="index.html" />
    </sync>
    <wrapper wrapped="img" :src="user.AvatarUrl|string" :wight="user.Wight|uint" />
</div>
```

Template file index.html :
``` Html
<div>
    <for itr_count=5>
        <h1>
            Hello from template!!!
        </h1>
    </for>
</div>
```

Result :
``` Html
<div>
   <h1>Hello World!</h1>
   <h1>Hello World!</h1>
   <h1>Hello World!</h1>
   <div>
      <h1>Hello from template!!!</h1>
      <h1>Hello from template!!!</h1>
      <h1>Hello from template!!!</h1>
      <h1>Hello from template!!!</h1>
      <h1>Hello from template!!!</h1>
   </div>
   <img src="test/img/bitch" wight="1"/>
</div>
```

