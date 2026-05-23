# whoistrader
CS2 trader profiler, built by players for players. Nobody likes scammers and reversers.

In a nutshell, it takes multiple APIs from CS2 marketplaces, reversal DBs, and Steam itself and aggregates them into one fat API response that represents the "profile" of a user. 

This project itself is not meant to explicitly qualify users as trustworthy or not; it instead provides information for potential buyers/sellers to access to make a decision.

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
