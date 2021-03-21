<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 512 512" style="height: 100px;">
  <g>
    <path d="M61 140c-8 0-14-6-14-14V61c0-8 6-14 14-14h65a14 14 0 110 28H75v51c0 8-6 14-14 14z" fill="#ef71a8"/>
    <path d="M451 465h-65a14 14 0 110-28h51v-51a14 14 0 1128 0v65c0 8-6 14-14 14z" fill="#ef71a8"/>
    <path d="M498 512H14c-8 0-14-6-14-14V14C0 6 6 0 14 0h484c8 0 14 6 14 14v484c0 8-6 14-14 14zM28 484h456V28H28v456z" fill="#ef71a8"/>
  </g>
  <path d="M369 266c8-21 12-44 12-65v-1c0-4 2-7 5-9 10-9 26-30 26-30l-40-6c-2 0-12-11-14-12a68 68 0 00-66-14c-22 7-39 25-44 47-2 8-3 17-2 25v2a3 3 0 01-2 1 178 178 0 01-122-65c-2-2-5-2-6 0a68 68 0 0011 82c-5-1-11-3-15-6-3-1-6 1-6 3 0 28 16 52 40 63a70 70 0 01-13-1c-3-1-5 2-4 4 8 24 29 42 53 46-20 14-44 21-69 21h-8c-3 0-5 2-5 4-1 3 0 5 2 6 29 17 61 25 94 25a176 176 0 00138-61" fill="#71cad1"/>
  <path d="M196 401c-36 0-70-10-101-28-7-4-11-13-9-21 2-9 10-15 19-15h8c12 0 24-2 35-5a82 82 0 01-33-44c-1-3-1-8 1-11a82 82 0 01-24-59c0-5 3-10 7-14a82 82 0 015-72 18 18 0 0129-2 164 164 0 0099 58l2-15a80 80 0 0179-62 82 82 0 0157 23l7 7 38 6a14 14 0 019 23c-2 2-17 21-29 31 0 23-4 47-13 70a14 14 0 11-26-10c7-20 11-40 11-60v-1c0-8 4-15 10-20l9-9-15-2c-6-1-11-5-21-14l-1-2a53 53 0 00-53-11c-17 6-30 20-34 37-2 7-2 14-2 20a17 17 0 01-18 19 191 191 0 01-120-57 54 54 0 0015 50 14 14 0 01-14 24c4 14 15 26 29 33a14 14 0 01-2 26c8 12 21 20 35 22a14 14 0 015 26c-15 10-32 17-50 21a174 174 0 00129-6c21-10 39-23 54-41a14 14 0 1121 19 190 190 0 01-148 66zm49-211zm127-21zm-261-21zm262-8z" fill="#ef71a8"/>
</svg>

# TweetSVG

Generate an SVG for a given Tweet ID

![](https://tweet.cluster.fun/1363048182020792325)

Available at https://tweet.cluster.fun/

## Features

* Provide the ID of a tweet and have it render as an SVG for use in an `<img>` tag.
* No JavaScript required

## Building from source

With Docker:

```sh
make docker-build
```

Standalone:

```sh
make build
```

## Contributing

If you find a bug or have an idea for a new feature please [raise an issue](issues/new) to discuss it.

Pull requests are welcomed but please try and follow similar code style as the rest of the project and ensure all tests and code checkers are passing.

Thank you 💛

## License

See [LICENSE](LICENSE)
