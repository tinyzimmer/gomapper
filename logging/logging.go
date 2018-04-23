/**
    This file is part of gomapper.

    Gomapper is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    Gomapper is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with gomapper.  If not, see <http://www.gnu.org/licenses/>.
**/

package logging

import (
	"fmt"
	"log"

	"github.com/tinyzimmer/gomapper/formats"
)

func LogError(msg string) {
	line := fmt.Sprintf("%s %s", formats.ColorRed("ERROR"), msg)
	log.Println(line)
}

func LogInfo(msg string) {
	line := fmt.Sprintf("%s %s", formats.ColorGreen("INFO"), msg)
	log.Println(line)
}

func LogWarn(msg string) {
	line := fmt.Sprintf("%s %s", formats.ColorYellow("WARNING"), msg)
	log.Println(line)
}

func LogDebug(msg string) {
	line := fmt.Sprintf("%s %s", formats.ColorBlue("DEBUG"), msg)
	log.Println(line)
}
