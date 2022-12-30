# nythead - print headlines from the New York Times

## usage

```
usage of nythead:
  -s string
    	headline type (arts, health, sports, science, technology, u.s., world) (default "u.s.")
```

## examples

```
$ nythead
Key Takeaways From Trump’s Tax Returns
Trump Tax Returns Released by House Democrats
Reported Sexual Abuse at California Prep School Won’t Be Prosecuted
Biden Signs Government Funding Bill, Preventing Shutdown
Southwest Is California’s ‘Unofficial Airline.’ The Meltdown Has Residents Anxious.
```

```
$ nythead -s arts
Giancarlo Esposito Plays Other People So He Can Know Himself
These Young Musicians Made an Album. Now It’s Nominated for a Grammy.
The Complex History Behind a Vienna Philharmonic Tradition
5 Things to Do This Weekend
Ian Tyson, Revered Canadian Folk Singer, Dies at 89
```

This program looks for the NYT API key in the environment variable ```NYTAPIKEY```, and fails if is not set.