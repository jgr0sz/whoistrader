# whoistrader
CS2 trader profiler, made by a player who isn't fond of the amount of scammers/reversers in the community.

In a nutshell, it takes various sources of trader information (reversal databases, marketplace APIs, Steam's web API) and builds a profile around them, displaying relevant information such as reversals, Steam-related bans, and statistics/standings on major marketplaces.

## Usage 
- In the root folder, create an `.env` file to house API keys for the platforms that require them.
- For the time being: the structure is as follows: 
<pre>
CSFLOAT_API_KEY=API_KEY
STEAM_API_KEY=API_KEY 
</pre>


You can generate a CSFloat API key at https://csfloat.com/profile > Developers > New Key.

You can generate a Steam Web API key at https://steamcommunity.com/dev/apikey. Enter anything you'd like in the domain name field and click "Register".

## Screenshots (as of 5/23)
<img width="1081" height="796" alt="whoistrader_sample" src="https://github.com/user-attachments/assets/278a6d69-ce2b-44c5-8215-81637841778d" />
(Sample output from two different steamIDs)

## Contributing
Any and all pull requests adding features or new platforms are welcome. [This](endpoints/csfloat.go) is a solid example of the modular design involved, the main components being a struct for JSON objects and `Fetch()/Name()` functions to satisfy the `Endpoint` interface.

## Plans
- Integrating as many reliable APIs and data sources as possible
- [ ] Exposing service responses over an HTTP server
- [ ] Discord bot integration (in a separate repo)
- [ ] Adding analysis capabilities based on trader data (Bayesian average score or similar)
