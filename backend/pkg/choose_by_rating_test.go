package pkg

import (
	"fmt"
	"math"
	"testing"
)

func TestChooseByRating(t *testing.T) {
	pow := 0.01
	user1Rating := 100.0
	user2Rating := 99.0
	user3Rating := 1.0
	new1Rating := math.Pow(user1Rating, pow)
	new2Rating := math.Pow(user2Rating, pow)
	new3Rating := math.Pow(user3Rating, pow)
	fmt.Println("New rating of first", new1Rating)
	fmt.Println("New rating of second", new2Rating)
	fmt.Println("New rating of third", new3Rating)
	fmt.Println("Probability with func for first", new1Rating/(new1Rating+new2Rating+new3Rating))
	fmt.Println("Probability with func for second", new2Rating/(new1Rating+new2Rating+new3Rating))
	fmt.Println("Probability with func for third", new3Rating/(new1Rating+new2Rating+new3Rating))

	//Result with pow=0.01:
	//New rating of first  1.0471285480508996
	//New rating of second 1.047023313403309
	//New rating of third  1

	//Probability with func for first  0.3384218341367265
	//Probability with func for second 0.3383878233149237
	//Probability with func for third  0.3231903425483499
}
