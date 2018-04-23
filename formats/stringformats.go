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

package formats

import (
	"fmt"
	"strconv"
)

func FloatToString(value float64) string {
	return strconv.FormatFloat(value, 'g', 4, 64)
}

func ColorRed(value string) string {
	return fmt.Sprintf("\033[0;31m%s\033[0m", value)
}

func ColorGreen(value string) string {
	return fmt.Sprintf("\033[0;32m%s\033[0m", value)
}

func ColorYellow(value string) string {
	return fmt.Sprintf("\033[0;33m%s\033[0m", value)
}

func ColorBlue(value string) string {
	return fmt.Sprintf("\033[0;34m%s\033[0m", value)
}
