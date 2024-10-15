package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/astoyanov87/live-score-service/models"
	"github.com/chromedp/chromedp"
)

func FetchLiveScore(matchId string) (models.Liveresult, error) {
	// Create a context for chromedp
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var matchSectionContent string
	fmt.Println("Scraping page : https://www.wst.tv/match-centre/" + matchId)
	chromedp.Run(ctx,
		chromedp.Navigate("https://www.wst.tv/match-centre/"+matchId),
		// Wait for the match data to be loaded
		chromedp.WaitVisible(`section.match-hero`),
		// Scrape the HTML content of the matches section
		chromedp.OuterHTML(`section.match-hero`, &matchSectionContent),
	)
	var liveScore models.Liveresult
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(matchSectionContent))
	if err != nil {
		panic(err)
	}

	firstSection := dom.Find("section.match-player-card").First()
	resScoreFirst := strings.TrimSpace(firstSection.Find("div.points p.text-lg").Text())
	resFramesFirst := strings.TrimSpace(firstSection.Find("div.frames p.text-lg").Text())

	resFramesFirstInt, err := strconv.Atoi(resFramesFirst)
	resScoreFirstInt, err := strconv.Atoi(resScoreFirst)

	breakFirst := strings.TrimSpace(firstSection.Find("div.break-wrapper p.text-base").Text())
	breakFirstInt, err := strconv.Atoi(breakFirst)

	fmt.Println("Current break 1: ")
	fmt.Println(breakFirst)

	secondSection := dom.Find("section.match-player-card").Last()
	resScoreSecond := strings.TrimSpace(secondSection.Find("div.points p.text-lg").Text())
	resScoreSecondInt, err := strconv.Atoi(resScoreSecond)

	resFramesSecond := strings.TrimSpace(secondSection.Find("div.frames p.text-lg").Text())
	resFramesSecondInt, err := strconv.Atoi(resFramesSecond)

	breakSecond := strings.TrimSpace(secondSection.Find("div.break-wrapper p.text-base").Text())
	breakSecondInt, err := strconv.Atoi(breakSecond)

	if err != nil {
		fmt.Println("Error converting string value to integer")
		return liveScore, err
	}
	// prepare return object
	liveScore.HomePlayerFrames = resFramesFirstInt
	liveScore.HomeplayerPointsInCurrentFrame = resScoreFirstInt
	liveScore.HomePlayerCurrentBreak = breakFirstInt
	liveScore.AwayPlayerFrames = resFramesSecondInt
	liveScore.AwayPlayerPointsInCurrentFrame = resScoreSecondInt
	liveScore.AwayPlayerCurrentBreak = breakSecondInt

	fmt.Printf(" Home player frames: %d \n # Score in current frame:%d", liveScore.HomePlayerFrames, liveScore.HomeplayerPointsInCurrentFrame)
	fmt.Printf(" Home player current break: %d \n", liveScore.HomePlayerCurrentBreak)
	fmt.Printf(" Away player frames: %d \n # Score in current frame:%d", liveScore.AwayPlayerFrames, liveScore.AwayPlayerPointsInCurrentFrame)
	fmt.Printf(" Away player current break: %d \n", liveScore.AwayPlayerCurrentBreak)

	return liveScore, nil
}
