/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package ui

var Mutation = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
    <link href="style/style.css" rel="stylesheet" type="text/css" />
    <link href="style/go4api.css" rel="stylesheet" type="text/css"/>
    <script type="text/javascript" src="js/go4api.js"></script>
    <script type="text/javascript" src="js/reslts.js"></script>
    <script type="text/javascript" src="js/stats.js"></script>
    <script type="text/javascript" src="js/executed.js"></script>
    <script type="text/javascript" src="js/notexecuted.js"></script>
    <script type="text/javascript" src="js/mutationstats.js"></script>
  <title>Go4Api Reports</title>
</head>
<body>
  <div class="container">
      <div class="head">
          <a href="https://github.com/zpsean/go4api" target="blank_" title="Go4Api Home Page"><img alt="Go4Api" src="style/logo.png"/></a>
      </div>
      <div class="main">
          <div class="skeleton">
              <div class="content">
                  <div class="sous-menu">
                      <div class="item "><a href="index.html">Overview</a></div>
                      <div class="item "><a id="graphic_link" href="graphic.html">Graphic</a></div>
                      <div class="item "><a id="details_link" href="details.html">Details</a></div>
                      <div class="item "><a id="mindex_link" href="mindex.html">MutationOverview</a></div>
                      <div class="item selected"><a id="mutation_link" href="mutation.html">Mutation</a></div>
                      <div class="item "><a id="fuzz_link" href="fuzz.html">Fuzz</a></div>

                      <script type="text/javascript">
                        // var timestamp = 1523957748602;
                        // var runStartHumanDate = moment(timestamp).format("YYYY-MM-DD HH:mm:ss Z");
                        document.writeln("<p class='sim_desc' title='" +"Started at 2018-xx-xx, duration : 10 seconds' data-content=''>");
                        document.writeln("<b>" + "Started at 2018-xx-xx, duration : 10 seconds </b>");
                        document.writeln("</p>");
                      </script>
                  </div>


                  <div class="content-in">
                    <div>
                        <select>
                          <option value ="Priority">Priority</option>
                          <option value ="Case ID">Case ID</option>
                          <option value="ParentTestCase">ParentTestCase</option>
                          <option value="Status">Status</option>
                        </select>

                        <input type="text" size="50" name="search_text" value="Please enter search text here">
                        <button type="button">Search!</button>
                    </div>

                    <h1><span>> </span>Overview Information</h1>
                    <div class="article">


                      <div class="statistics extensible-geant collapsed">
                          

                          <table id="container_statistics_head" class="statistics-in extensible-geant">
                              <thead>
                                  <tr>
                                      <th id="col-1" class="header sortable"><span>#</span></th>
                                      <th id="col-2" class="header sortable"><span>HttpUrl</span></th>
                                      <th id="col-2" class="header sortable"><span>HttpMethod</span></th>
                                      <th id="col-3" class="header sortable"><span>MutationPart</span></th>
                                      <th id="col-4" class="header sortable"><span>MutationCategory</span></th>
                                      <th id="col-4" class="header sortable"><span>MutationRule</span></th>
                                      <th id="col-4" class="header sortable"><span>HttpStatus</span></th>
                                      <th id="col-5" class="header sortable"><span>TestStatus</span></th>
                                      <th id="col-6" class="header sortable"><span>Count</span></th>
                                      <th id="col-7" class="header sortable"><span>MutationMessage</span></th>
                                  </tr>
                              </thead>
                              <tbody></tbody>
                          </table>

                          <div class="scrollable">
                              <table id="container_statistics_body" class="statistics-in extensible-geant">
                                  <tbody></tbody>
                              </table>


                              <script type="text/javascript">
                                for (var i = 0;i < tcResults.length; i++)
                                  {
                                    var newTr = container_statistics_body.insertRow();
                                    
                                    var newTd0 = newTr.insertCell();
                                    var newTd1 = newTr.insertCell();
                                    var newTd2 = newTr.insertCell();
                                    var newTd3 = newTr.insertCell();
                                    var newTd4 = newTr.insertCell();
                                    var newTd5 = newTr.insertCell();
                                    var newTd6 = newTr.insertCell();
                                    var newTd7 = newTr.insertCell();
                                    var newTd8 = newTr.insertCell();
                                    var newTd9 = newTr.insertCell();
                             
                                    newTd0.innerText = i;
                                    newTd1.innerText = tcResults[i].Path;
                                    newTd2.innerText = tcResults[i].Method;
                                    newTd3.innerText = tcResults[i].MutationArea;
                                    newTd4.innerText = tcResults[i].MutationCategory;
                                    newTd5.innerText = tcResults[i].MutationRule;
                                    newTd6.innerText = tcResults[i].ActualStatusCode;
                                    newTd7.innerText = tcResults[i].TestResult;
                                    newTd8.innerText = 1;
                                    newTd9.innerText = JSON.stringify(tcResults[i].MutationInfo, null, 4);
                                  }
                              </script>

                          </div>
                      </div>
                    </div>

                  </div>
              </div>
          </div>
      </div>
  </div>
  <div class="foot">
      <a href="https://github.com/zpsean/go4api" title="Go4Api Home Page"><img alt="Go4Api" src="style/logosmall.png"/></a>
  </div>
</body>
</html>
`
