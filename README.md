# CalHacks4-Backend
This project was built for CalHacks 4.0.  This repository is holds the backend
server code for the app.

Here are the links to the
[FrontEnd](https://github.com/ayush-patel/CalHacks4-iOS) and [Uploads
Scripts](https://github.com/SachitShroff/CalHacks4-processing)

## Inspiration
The hard to use, disorganized, and disjoint methods of presenting course
information and lectures at our schools (UCLA and Cal) inspired us to create a
unified interface.

## What It Does
Professors and students can create accounts. Professors can upload webcasts of
their lectures as well as any supplementary materials. Students can register for
classes, where they'll be able to watch all the lectures, search by topic, view
and add supplementary materials (like worksheets or notes) andother helpful
resources.

## How I Built It
The front end is built in Xcode using Swift (as an iPad app). The front end
pings a REST API hosted on a Go web server (which also hosts the database). The
web server authenticates users and handles file uploads and post requests.
Webcasts are stored in Google Cloud Storage to allow for easy access and
integration of the Google Cloud Speech API. The python script we created
processes and annotates the videos (using speech-to-text from Google Cloud
Speech API, Keyword Analysis from Azure Cognitive Services, and a custom program
to find relevant online results), and sends the relevant info to the web server.

## Challenges I ran Into
Designing a clean UI, API limitations (with beta Speech API client library for
instance) and quotas, creating a robust web-server, error handling (especially
in communications between different components of program) and extracting videos
and audio in the correct formats.

## Accomplishments that I'm proud of
Strung together output from a variety of useful APIs to generate useful
information. Gaining a deep understanding of Google Cloud Speech API and
discovering errors with documentation that allowed us to provide Google with
useful feedback. Handled errors across a variety of usecases (actually robust
program).

## What I learned
How to deal with processing/converting audio and video efficiently to prevent
bandwidth choking. Knowledge about the interaction between databases/services.

## What next for Intellicast
Creating an interactive interface that allows students to contribute and
interact, and intelligent annotation/parsing of a variety of related materials
from professors.
