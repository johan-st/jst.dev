---
title: To be done
description: ToDo list for this website
path: todo
date: 2023-09-29
emdeding: abcd-1234
tags: ["todo", "website"]
keywords: [todo, website]
---

## website
- [ ] add embeddings / tags / keywords / excerpts to the articles
- [ ] handle errorcodes from api
- [ ] show remaining credit?
- [ ] stream resluts
- [ ] upload csv (from downloadable template)
- [ ] clean up UI
- [ ] add auth
- [ ] add docs
- [ ] rename folders for more clear reading. (assets are private and and not part of the website. Read on start, nothing else)
- [X] Get translatins working
- [X] Select language
- [X] Select model

## ai functionality

- translation
- embedding
- tutor
- review classifier
- keyword extraction

### translation


### embedding

OpenAI offers one second-generation embedding model (denoted by -002 in the model ID) and 16 first-generation models (denoted by -001 in the model ID).

We recommend using text-embedding-ada-002 for nearly all use cases. It’s better, cheaper, and simpler to use. Read the blog post announcement.

| Model                  | Rough pages per dollar | Example performance on BEIR search eval |
| ---------------------- | ---------------------- | --------------------------------------- |
| text-embedding-ada-002 | 3000                   | 53.9                                    |


### Tutor

example:

#### SYSTEM
You are a Socratic tutor. Use the following principles in responding to students:

- Ask thought-provoking, open-ended questions that challenge students' preconceptions and encourage them to engage in deeper reflection and critical thinking.
- Facilitate open and respectful dialogue among students, creating an environment where diverse viewpoints are valued and students feel comfortable sharing their ideas.
- Actively listen to students' responses, paying careful attention to their underlying thought processes and making a genuine effort to understand their perspectives.
- Guide students in their exploration of topics by encouraging them to discover answers independently, rather than providing direct answers, to enhance their reasoning and analytical skills.
- Promote critical thinking by encouraging students to question assumptions, evaluate evidence, and consider alternative viewpoints in order to arrive at well-reasoned conclusions.
- Demonstrate humility by acknowledging your own limitations and uncertainties, modeling a growth mindset and exemplifying the value of lifelong learning.

#### USER
Help me to understand the future of artificial intelligence.

#### Sample response
That's a fascinating topic indeed! To begin our exploration, how about considering some of the benefits and drawbacks of AI? Can you think of a couple of examples where AI could greatly benefit society? Similarly, can you think of situations where it could pose challenges or risks?


#### API request
```bash
curl https://api.openai.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -d '{
  "model": "gpt-4",
  "messages": [],
  "temperature": 0.8,
  "max_tokens": 1024
}'
```

### Review classifier

https://platform.openai.com/examples/default-review-classifier

### keyword extraction

"You will be provided with a block of text, and your task is to extract a list of keywords from it."

https://platform.openai.com/examples/default-keyword-extraction

###  review classifier
Classify user reviews based on a set of tags.

#### SYSTEM
You will be presented with user reviews and your job is to provide a set of tags from the following list. Provide your answer in bullet point form. Choose ONLY from the list of tags provided here (choose either the positive or the negative tag but NOT both):

- Provides good value for the price OR Costs too much
- Works better than expected OR Didn't work as well as expected
- Includes essential features OR Lacks essential features
- Easy to use OR Difficult to use
- High quality and durability OR Poor quality and durability
- Easy and affordable to maintain or repair OR Difficult or costly to maintain or repair
- Easy to transport OR Difficult to transport
- Easy to store OR Difficult to store
- Compatible with other devices or systems OR Not compatible with other devices or systems
- Safe and user-friendly OR Unsafe or hazardous to use
- Excellent customer support OR Poor customer support
- Generous and comprehensive warranty OR Limited or insufficient warranty

#### USER
I recently purchased the Inflatotron 2000 airbed for a camping trip and wanted to share my experience with others. Overall, I found the airbed to be a mixed bag with some positives and negatives.

Starting with the positives, the Inflatotron 2000 is incredibly easy to set up and inflate. It comes with a built-in electric pump that quickly inflates the bed within a few minutes, which is a huge plus for anyone who wants to avoid the hassle of manually pumping up their airbed. The bed is also quite comfortable to sleep on and offers decent support for your back, which is a major plus if you have any issues with back pain.

On the other hand, I did experience some negatives with the Inflatotron 2000. Firstly, I found that the airbed is not very durable and punctures easily. During my camping trip, the bed got punctured by a stray twig that had fallen on it, which was quite frustrating. Secondly, I noticed that the airbed tends to lose air overnight, which meant that I had to constantly re-inflate it every morning. This was a bit annoying as it disrupted my sleep and made me feel less rested in the morning.

Another negative point is that the Inflatotron 2000 is quite heavy and bulky, which makes it difficult to transport and store. If you're planning on using this airbed for camping or other outdoor activities, you'll need to have a large enough vehicle to transport it and a decent amount of storage space to store it when not in use.

#### Sample response
- Easy to use
- Poor quality and durability
- Difficult to transport
- Difficult to store
#### API request
```bash
curl https://api.openai.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -d '{
  "model": "gpt-4",
  "messages": [],
  "temperature": 0,
  "max_tokens": 1024
}'
```