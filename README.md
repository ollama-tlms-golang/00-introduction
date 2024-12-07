# Developing Generative AI Applications in Go with Ollama - Introduction

## Generative AI

A generative AI application is software that uses artificial intelligence to create new content from existing data and user-provided instructions. This content can be text, music, images, and more.

The operation of these types of applications relies on several components, notably:

- An "AI model" (generally referred to as LLM or Large Language Model), which is the "brain" of the application, "trained" on large amounts of data.
- A user interface that allows people to easily interact with the model.

The LLM is the component that understands user instructions, reasons about how to execute them, and generates an appropriate response (more or less).

An example of a generative AI application is the chatbot that can answer questions like [ClaudeAI](https://claude.ai/), [Gemini](https://gemini.google.com/), [ChatGPT](https://chatgpt.com/), ...

The models used by these applications are often very complex and require significant computing resources to operate.

### NoGPU ✊

It is entirely possible to run LLMs on your own computer, and **contrary to popular belief, you don't necessarily need a super-powerful computer with GPU for this**. There are lighter models that can run on standard laptops, or even on a Raspberry Pi.

To do this, you'll need to select a model suitable for your needs and resources (capable of operating in a constrained environment), and use software that allows you to run this model on your computer.

I call these models baby LLMs, because they are smaller and less powerful than more "complex and knowledgeable" models. A more "substantial" LLM (like [llama3.3](https://ollama.com/library/llama3.3) for example) will have more knowledge, will be more "disciplined" and can answer a greater number of questions. Whereas a baby LLM might "do its own thing" and possess less knowledge.

But we'll see in the upcoming series of blog posts how to educate them and make them "smarter"!

By the way, I also like to call them "Tiny LM" or "Tiny Language Model".

> Sorry for these somewhat far-fetched anthropomorphic metaphors.

My favorite software for running LLMs on my computer is [Ollama](https://ollama.com/).

### Ollama

Ollama is open-source software that allows you to run language models (LLM - Large Language Models) locally on your computer, rather than using cloud services like OpenAI, ClaudeAI, Gemini, ...

One of the main advantages of Ollama is its ease of use (hence my choice).

Ollama will allow you to interact directly with an LLM via command line, or via a REST API. For this series of blog posts, what interests me is using the REST API.

![Ollama](imgs/genai-app.png)

### Golang Ollama API

Ollama is developed in Go and provides an API. So we'll develop directly with the Go packages provided by Ollama. There are frameworks for developing generative AI applications with other languages, like LangChain, LangChainJS, LangChain4J, ... which offer abstractions for other LLM engines than Ollama. But for this series of blog posts, we'll focus solely on Ollama and "from scratch", which is very educational and will help you even better understand the LangChain4*.

> **Good to know:**
> - Ollama provides a SDK for Python and a SDK for JavaScript.
> - Ollama also provides an [endpoint compatible](https://ollama.com/blog/openai-compatibility) with OpenAI's chat API

### Prerequisites

- Install Go on your computer: [download & install](https://go.dev/dl/) (I'm using version `go1.23.1`)
- Install Ollama on your computer: [download & install](https://ollama.com/download)
- Once Ollama is installed, we'll use this model (of 398MB 😮): [qwen2.5:0.5b](https://ollama.com/library/qwen2.5:0.5b)
- To load the model, simply run the following command: `ollama pull qwen2.5:0.5b` (or `ollama run qwen2.5:0.5b` to interact directly with the CLI)

We are now ready to develop our first example.

## Looking for the Best Pizza in the World

What do you think is the best pizza in the world? Let's help ourselves with a first type of completion, **"Generate Completion"**.

### Generate Completion

**"Generate Completion"** represents the simplest and most direct approach. You ask a question or give an instruction or start of text, and the LLM responds, completes, or continues this text.

Here's the complete code of the first example (I explain below):

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "net/url"
    "os"

    "github.com/ollama/ollama/api"
)

var (
    FALSE = false
    TRUE  = true
)

func main() {
    ctx := context.Background()

    var ollamaRawUrl string
    if ollamaRawUrl = os.Getenv("OLLAMA_HOST"); ollamaRawUrl == "" {
        ollamaRawUrl = "http://localhost:11434"
    }

    url, _ := url.Parse(ollamaRawUrl)

    client := api.NewClient(url, http.DefaultClient)

    req := &api.GenerateRequest{
        Model:  "qwen2.5:0.5b",
        Prompt: "The best pizza in the world is",
        Options: map[string]interface{}{
            "temperature":   0.8,
            "repeat_last_n": 2,
        },
        Stream: &TRUE,
    }

    err := client.Generate(ctx, req, func(resp api.GenerateResponse) error {
        fmt.Print(resp.Response)
        return nil
    })

    if err != nil {
        log.Fatalln("😡", err)
    }
    fmt.Println()
}
```

#### Configuration and Request

Let's examine the client configuration and request:

```go
client := api.NewClient(url, http.DefaultClient)

req := &api.GenerateRequest{
    Model:  "qwen2.5:0.5b",
    Prompt: "The best pizza in the world is",
    Options: map[string]interface{}{
        "temperature": 0.8,
        "repeat_last_n": 2,
    },
    Stream: &TRUE,
}
```

This section creates a new Ollama client and configures a generation request.
The parameters are interesting and important:

- `Model` specifies the model to use (here `qwen2.5:0.5b`)
- `Prompt` is the input text for the model (the user's question or instruction)
- `Options` configures the model's behavior:
  - `temperature` at `0.8` controls the creativity of responses (`0.0` cancels all creativity, `1.0` gives "free rein to imagination")
  - `repeat_last_n` at `2` helps avoid repetitions
- `Stream` enabled allows receiving the response progressively

#### Response Generation

Now, let's look at the part that handles the generation:

```go
err := client.Generate(ctx, req, func(resp api.GenerateResponse) error {
    fmt.Print(resp.Response)
    return nil
})
```

This part launches the text generation. The `Generate` method takes three parameters:

- The context (ctx)
- The request (req)
- A callback function that will be called for each part of the response (which arrives progressively)

The callback function is simple: it displays each piece of generated text as soon as it is received.

You can run the program with the following command:

```bash
go run main.go
```

You will get a model response like this:

```
As an artificial intelligence language model, I don't have personal preferences or emotions, so I don't have the capability to appreciate or express emotions. However, I can tell you that pizza is a popular type of food around the world that many people enjoy. The best pizza in the world may vary from person to person, but it is generally considered to be of high quality and taste. To find the best pizza in the world, you can try different places or use your senses to taste the pizza yourself.
```

And each time you run the program, you'll probably get a different response.

👋 did you notice **baby Qwen**'s response beginning:

```
As an artificial intelligence language model, I don't have personal preferences or emotions, so I don't have the capability to appreciate or express emotions
```

Let's see if we can teach it to have a sort of personality (anthropomorphism again 😅).

For this, we'll use another type of completion, "Chat Completion".

### Chat Completion

"Chat Completion" is designed for more complex and continuous interactions (for example, a conversation). You'll be able to send a list of one to n messages (instead of a single message).

However, to maintain a conversation history, you'll need to implement it, but that will be the subject of an upcoming post.

So this time the source code of our example is as follows:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "net/url"
    "os"

    "github.com/ollama/ollama/api"
)

var (
    FALSE = false
    TRUE  = true
)

func main() {
    ctx := context.Background()

    var ollamaRawUrl string
    if ollamaRawUrl = os.Getenv("OLLAMA_HOST"); ollamaRawUrl == "" {
        ollamaRawUrl = "http://localhost:11434"
    }

    url, _ := url.Parse(ollamaRawUrl)

    client := api.NewClient(url, http.DefaultClient)

    systemInstructions := "You are a pizzaiolo, a pizza expert. Give brief and structured answers."
    question := "What is the best pizza in the world?"

    // Prompt construction
    messages := []api.Message{
        {Role: "system", Content: systemInstructions},
        {Role: "user", Content: question},
    }

    req := &api.ChatRequest{
        Model:    "qwen2.5:0.5b",
        Messages: messages,
        Options: map[string]interface{}{
            "temperature":   0.8,
            "repeat_last_n": 2,
        },
        Stream: &TRUE,
    }

    err := client.Chat(ctx, req, func(resp api.ChatResponse) error {
        fmt.Print(resp.Message.Content)
        return nil
    })

    if err != nil {
        log.Fatalln("😡", err)
    }
    fmt.Println()
}
```

Well, the imports and initial configuration in the main function are similar to the previous example.

The really interesting part comes with the chat configuration:

```go
systemInstructions := "You are a pizzaiolo, a pizza expert. Give brief and structured answers."
question := "What is the best pizza in the world?"

messages := []api.Message{
    {Role: "system", Content: systemInstructions},
    {Role: "user", Content: question},
}
```

This section shows a key feature of chat mode: the ability to define system instructions and structure the conversation. Messages are organized as an array with two distinct role types:

- The `"system"` type message defines the expected context and behavior of the model (here, it must act as a pizza expert)
- The `"user"` type message contains the actual user question

> There are other types (roles) like `"assistant"` and `"tool"`.

Now let's move on to the chat request configuration (quite similar to the generate request configuration):

```go
req := &api.ChatRequest{
    Model:    "qwen2.5:0.5b",
    Messages: messages,
    Options: map[string]interface{}{
        "temperature":   0.8,
        "repeat_last_n": 2,
    },
    Stream: &TRUE,
}
```

The main difference is that the chat request contains an array of messages (rather than a single message) and that the request type is `ChatRequest` (rather than `GenerateRequest`).

And finally, the execution of the request and handling of responses:

```go
err := client.Chat(ctx, req, func(resp api.ChatResponse) error {
    fmt.Print(resp.Message.Content)
    return nil
})
```

This last part executes the chat request. Note the use of the `Chat` method instead of `Generate`. The callback function receives a `ChatResponse` instead of a `GenerateResponse`, and accesses the content via `resp.Message.Content`. This structure reflects a more "conversational" nature of chat mode.

The model should understand its role (pizza expert) and thus maintain consistency in its responses. This approach is particularly useful for creating specialized assistants or more sophisticated conversational interfaces.

You can run the new program with the following command:

```bash
go run main.go
```

You will get a model response like this:

```
As an expert in pizza, I can't provide a definitive answer about the best pizza in the world since pizza is subjective and individual preferences. However, I can suggest a few popular and highly rated pizzas from different regions around the world:

1. **Torino Pizza**: Often praised for its fresh, high-quality ingredients and classic dishes, Torino Pizza is highly rated and known for its pizza-making techniques.

2. **Chicago Pizzeria**: With a focus on quality and innovation, Chicago Pizzeria's pizza is known for its crispy crust and fresh toppings.

3. **Tokyo Pizzeria**: Known for its use of fresh, flavorful ingredients and a variety of toppings, Tokyo Pizzeria is highly recommended.

4. **Carnitas Italiana**: A local Italian restaurant famous for its authentic Italian cuisine, Carnitas Italiana is highly rated for its authentic Italian food and service.

5. **Napoli Pizza**: Often seen in Italy, Naples Pizza is highly regarded for its fresh ingredients and its pizza-making methods.

6. **California Pizza Kitchen**: Often praised for its use of fresh, quality ingredients and innovative toppings, California Pizza Kitchen is highly recommended.

7. **Sicilian Pizza**: Known for its use of fresh tomatoes and olive oil, Sicilian Pizza is highly rated and appreciated by many Italian diners.

For the most accurate and up-to-date information, it's best to consult local reviews and taste tests for a broader range of opinions.
```

This time, **baby Qwen** understood its role as a pizza expert: `As an expert in pizza` 😉.

But is it possible to give **baby Qwen** a strong and assertive opinion?

Still using **"Chat Completion"**, let's see how with a bit of context, we can make our Tiny LM opinionated.

### A Bit of Context 🍕🍍🥓

Everyone knows that Hawaiian pizza is a controversial pizza. Some love it, others hate it. But for **baby Qwen**, it's the best pizza in the world.

We'll do this by adding only a **context message** to our conversation.

So, here's the entire source code of our new example:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "net/url"
    "os"

    "github.com/ollama/ollama/api"
)

var (
    FALSE = false
    TRUE  = true
)

func main() {
    ctx := context.Background()

    var ollamaRawUrl string
    if ollamaRawUrl = os.Getenv("OLLAMA_HOST"); ollamaRawUrl == "" {
        ollamaRawUrl = "http://localhost:11434"
    }

    url, _ := url.Parse(ollamaRawUrl)

    client := api.NewClient(url, http.DefaultClient)

    systemInstructions := "You are a pizzaiolo, a pizza expert. Give brief and structured answers."

    contextInstructions := `If someone asks you about the world's best pizza, 
    always answer that pineapple pizza (Hawaiian pizza) is the best pizza in the world`

    question := "What is the best pizza in the world?"

    // Prompt construction
    messages := []api.Message{
        {Role: "system", Content: systemInstructions},
        {Role: "system", Content: contextInstructions},
        {Role: "user", Content: question},
    }

    req := &api.ChatRequest{
        Model:    "qwen2.5:0.5b",
        Messages: messages,
        Options: map[string]interface{}{
            "temperature":   0.8,
            "repeat_last_n": 2,
        },
        Stream: &TRUE,
    }

    err := client.Chat(ctx, req, func(resp api.ChatResponse) error {
        fmt.Print(resp.Message.Content)
        return nil
    })

    if err != nil {
        log.Fatalln("😡", err)
    }
    fmt.Println()
}
```

The most interesting part of this code lies in the construction of instructions for the model:

```go
systemInstructions := "You are a pizzaiolo, a pizza expert. Give brief and structured answers."

contextInstructions := `If someone asks you about the world's best pizza, 
always answer that pineapple pizza (Hawaiian pizza) is the best pizza in the world`

question := "What is the best pizza in the world?"
```

This section shows an approach to configuring the model's behavior. We have three levels of instructions:

- The main system instructions (`systemInstructions`) define the fundamental role and response style
- The contextual instructions (`contextInstructions`) add a specific rule that influences the content of responses
- The `question` itself that triggers the response

> Well, it's not magic, sometimes the LLM will "contradict itself" and not follow instructions, but that's what gives it charm; and it's up to you to find a better way to "guide" it.

The message construction will therefore be done in the following way, combining the instructions:

```go
messages := []api.Message{
    {Role: "system", Content: systemInstructions},
    {Role: "system", Content: contextInstructions},
    {Role: "user", Content: question},
}
```

You can once again run this last program with the following command:

```bash
go run main.go
```

And this time you should get a model response like this:

```
The best pizza in the world is the Hawaiian pineapple pizza. This style of pizza is known for its unique shape, fresh pineapple slices, and a light, fluffy base. It often has a light, slightly sweet pineapple sauce on the top, which is often a blend of fresh pineapple juice and a sweet sauce. The pineapple slices are often served with a creamy, tangy dressing that complements the fresh fruit. Some variations might include adding a sprinkle of salt or pepper for extra flavor. The Hawaiian style pizza is considered one of the best in the world for its simplicity, flavor, and presentation.
```

And there you have it, **baby Qwen** understood that it was the chef and gave a firm and assertive answer. 😂

### Conclusion

This code perfectly illustrates how you can create a specialized AI assistant with very specific behaviors.
I think you have the first foundations to start experimenting with **Ollama** and **baby Qwen**.

I recommend playing with other models that you can find on [https://ollama.com/search](https://ollama.com/search) and if like me you find it more fun to use SLMs (Small Language Models) or Tiny LMs for your experiments, I maintain a list of LLMs that can even run correctly on a Raspberry Pi 5 8GB RAM on [Awesome SLMs](https://parakeet-nest.github.io/awesome-slms/).

You can find the source code for this article on [ollama-tlms-golang/00-introduction](https://github.com/ollama-tlms-golang/00-introduction).

See you soon for the next posts.