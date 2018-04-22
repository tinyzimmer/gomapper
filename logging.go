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

package main

import (
	"fmt"
	"log"
)

func logError(msg string) {
	line := fmt.Sprintf("\033[0;31mERROR:\033[0m %s", msg)
	log.Println(line)
}

func logInfo(msg string) {
	line := fmt.Sprintf("\033[0;32mINFO:\033[0m %s", msg)
	log.Println(line)
}

func logWarn(msg string) {
	line := fmt.Sprintf("\033[0;33mWARNING:\033[0m %s", msg)
	log.Println(line)
}

func logDebug(msg string) {
	line := fmt.Sprintf("\033[0;34mDEBUG:\033[0m %s", msg)
	log.Println(line)
}
