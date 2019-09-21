<p align="center">
<a href="https://imgbb.com/"><img src="https://i.ibb.co/4Swf8gL/bird-singing2.jpg" alt="bird-singing2" border="0"></a><br />
</p>
<p align="center">
<br /><br /><br />

Simple CLI server for all your blabbing needs.

It packs two simple but powerful functionalities:

  1. You can give the server text, and it will process and learn its ngrams (n configurable when starting the server).
  2. You can ask the server to generate random text based on the ngrams it's learnt from the input text it's been fed since starting.

```
Usage of ./blabber:
  -maxWords uint
    	max number of words on output when using generate (default 100)
  -ngram uint
    	size of the ngrams to be processed when using learn (default 3)
  -port uint
    	port for the server to listen on (default 8080)
```

Example of use:

```
$ curl -X POST http://localhost:8080/learn -H "Content-Type: text/plain --data-binary @pride-prejudice.txt 
$ curl -X GET http://localhost:8080/generate
To think it more than commonly anxious to get round to the preference of one, and offended by the other as politely and more cheerfully. Their visit afforded was produced by the lady with whom she almost looked up to the stables. They were to set out with such a woman.
```

The possibilities are endless!
