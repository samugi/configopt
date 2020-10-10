package com.configopt;

import java.util.ArrayList;
import java.util.List;
import java.util.logging.Level;
import java.util.logging.Logger;

import com.configopt.Utils.Mode;
import org.apache.commons.cli.*;
import org.json.simple.JSONObject;

public class ConfigOpt {
    public static void main(String[] args) {

        String usageString = "java -cp target/configopt-1.0-jar-with-dependencies.jar com.configopt.ConfigOpt [options...] --configuration <arg>";
        Options options = new Options();

        Option configurationOption = Option.builder("c").longOpt("configuration").required(true)
                .desc("JSON configuration file path").hasArg().build();
        Option outputFileOption = Option.builder("o").longOpt("output").required(false).desc("Output file for report").hasArg()
                .build();
        Option debugLogLevelOption = Option.builder("v").longOpt("verbose").required(false).desc("Verbose logs")
                .build();
        Option interactiveModeOption = Option.builder("i").longOpt("interactive").required(false)
                .desc("Enables interactive mode").build();
        Option pathRoutingOnlyOption = Option.builder("pro").longOpt("pathroutingonly").required(false)
                .desc("Runs in path routing only mode. Use this is you have APICAST_PATH_ROUTING_ONLY=true").build();
        Option helpOption = Option.builder("h").longOpt("help").required(false)
                .desc("Show this help message").build();

        options.addOption(configurationOption).addOption(outputFileOption).addOption(debugLogLevelOption)
                .addOption(interactiveModeOption).addOption(helpOption).addOption(pathRoutingOnlyOption);

        CommandLineParser parser = new DefaultParser();
        HelpFormatter formatter = new HelpFormatter();
        CommandLine cmd = null;

        try {
            cmd = parser.parse(options, args);
        } catch (ParseException e) {
            System.out.println(e.getMessage());
            formatter.printHelp(usageString, options);
            System.exit(1);
        }
        if(cmd.hasOption(helpOption.getOpt())){
            formatter.printHelp(usageString, options);
            System.exit(1);
        }
        String inputFilePath = cmd.getOptionValue("configuration");
        OutputUtils.outputFile = cmd.getOptionValue("output");
        
        JSONObject JSONConfig = Utils.extractConfigJSONFromFile(inputFilePath);
        List<Service> services = Utils.createServicesFromJSONConfig(JSONConfig);
        if (cmd.hasOption(debugLogLevelOption.getOpt()))
            Logger.getLogger(Utils.LOG_TAG).setLevel(Level.ALL);
        else
            Logger.getLogger(Utils.LOG_TAG).setLevel(Level.SEVERE);
        if (cmd.hasOption(interactiveModeOption.getOpt()))
            Utils.mode = Mode.FIXINTERACTIVE;
        else
            Utils.mode = Mode.SCAN;
        if(OutputUtils.outputFile == null && Utils.mode == Mode.FIXINTERACTIVE)
            Logger.getLogger(Utils.LOG_TAG).log(Level.SEVERE, "No output file provided with fix-interactive mdoe enabled: will not print output configuration.");
      
        // Service service1 = new Service(12l, "example.org");
        // service1.addProductMappingRule(new MappingRuleSM("GET", "/foo/bar/test$", 12l, "proxy", 10l));
        // service1.addProductMappingRule(new MappingRuleSM("GET", "/foo/{bar}/test$", 12l, "proxy", 11l));
        // service1.addProductMappingRule(new MappingRuleSM("GET", "/foo/bar/test$", 12l, "proxy", 12l));
        // service1.addProductMappingRule(new MappingRuleSM("GET", "/test-tttt/v3.0/andf/accounts/{AccountId}$", 12l, null, 13l));
        // service1.addProductMappingRule(new MappingRuleSM("GET", "/fo", 12l, null, 14l));

        // Service service2 = new Service(13l, "example.org");
        // service2.addProductMappingRule(new MappingRuleSM("GET", "/foo/bar", 13l, null, 15l ));
        // service2.addProductMappingRule(new MappingRuleSM("GET", "/foo/bar/test", 13l, null, 16l ));
        // service2.addProductMappingRule(new MappingRuleSM("GET", "/whatever", 13l, null,17l));
        // service2.addProductMappingRule(new MappingRuleSM("GET", "/test-tttt/v3.0/$", 13l, null,18l));
        // service2.addProductMappingRule(new MappingRuleSM("GET", "/test-tttt/v3.0/andf/accounts/{AccountId}$", 13l, null,19l));
        // service2.addProductMappingRule(new MappingRuleSM("GET", "/test-tttt/v3.0/andf/accounts/{AccountId}$", 13l, null,20l));
        // service2.addProductMappingRule(new MappingRuleSM("GET", "/fo", 13l, null,21l));
        
        // List<Service> services = new ArrayList<>();
        // services.add(service1);
        // services.add(service2);

        APIcast apicast = APIcast.getAPIcast();
    //  if(cmd.hasOption(pathRoutingOption.getOpt()))
    //    apicast.setPathRoutingEnabled(true);
        if(cmd.hasOption(pathRoutingOnlyOption.getOpt()))
            apicast.setPathRoutingOnlyEnabled(true);
        for (Service service : services) {
            apicast.addService(service);
        }
        apicast.validateAllServices();

        if(Utils.mode == Mode.FIXINTERACTIVE){
            OutputUtils.rewriteConfig(apicast, JSONConfig);
        }
    }
}